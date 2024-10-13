package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	reportFilePath := os.Getenv("REPORT_FILE_PATH")
	namespace := os.Getenv("KUBERNETES_NAMESPACE")
	webhookURL := os.Getenv("WEBHOOK_URL")
	targetName := os.Getenv("TARGET_NAME")

	uniqueID := uuid.New().String()

	clientset, err := getKubernetesClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	waitForFile(reportFilePath)

	vulnerabilitiesFound := checkForVulnerabilities(reportFilePath)
	if vulnerabilitiesFound {
		err = handleVulnerabilities(clientset, targetName, namespace, reportFilePath, uniqueID, webhookURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error handling vulnerabilities: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("No vulnerabilities found.")
	}
}

func getKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating in-cluster config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %w", err)
	}
	return clientset, nil
}

func waitForFile(filePath string) {
	for {
		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("Report file found: %s\n", filePath)
			return
		}
		fmt.Println("Waiting for report file...")
		time.Sleep(5 * time.Second)
	}
}

func checkForVulnerabilities(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking report file: %v\n", err)
		return false
	}
	return fileInfo.Size() > 0
}

func handleVulnerabilities(clientset *kubernetes.Clientset, targetName, namespace, filePath, uniqueID, webhookURL string) error {
	err := createConfigMap(clientset, targetName, namespace, filePath, uniqueID)
	if err != nil {
		return fmt.Errorf("error creating ConfigMap: %w", err)
	}

	err = deployNginxPod(clientset, targetName, namespace, uniqueID)
	if err != nil {
		return fmt.Errorf("error deploying Nginx pod: %w", err)
	}

	err = createService(clientset, targetName, namespace, uniqueID)
	if err != nil {
		return fmt.Errorf("error creating Service: %w", err)
	}

	serviceURL := fmt.Sprintf("http://%s-dast-report-%s.%s.svc.cluster.local", targetName, uniqueID, namespace)
	message := fmt.Sprintf("Found vulnerabilities. To access the report, port-forward the service: %s or run kubectl port-forward svc/%s-dast-report-%s-service 8080:80 -n %s", serviceURL, targetName, uniqueID, namespace)
	sendWebhook(webhookURL, targetName, message, uniqueID)
	return nil
}

func createConfigMap(clientset *kubernetes.Clientset, targetName, namespace, filePath, uniqueID string) error {
	reportData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading report file: %w", err)
	}
	configMapName := fmt.Sprintf("%s-dast-report-configmap-%s", targetName, uniqueID)
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configMapName,
		},
		Data: map[string]string{
			"index.html": string(reportData),
		},
	}
	_, err = clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err = clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("error updating existing ConfigMap: %w", err)
			}
		} else {
			return fmt.Errorf("error creating ConfigMap: %w", err)
		}
	}
	return nil
}

func deployNginxPod(clientset *kubernetes.Clientset, targetName, namespace, uniqueID string) error {
	podName := fmt.Sprintf("%s-dast-report-%s", targetName, uniqueID)
	labels := map[string]string{
		"app": podName,
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   podName,
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:alpine",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "report-volume",
							MountPath: "/usr/share/nginx/html",
						},
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("10m"),
							corev1.ResourceMemory: resource.MustParse("10Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("20m"),
							corev1.ResourceMemory: resource.MustParse("20Mi"),
						},
					},
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "report-volume",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: fmt.Sprintf("%s-dast-report-configmap-%s", targetName, uniqueID),
							},
						},
					},
				},
			},
		},
	}
	_, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			err = clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				return fmt.Errorf("error deleting existing Pod: %w", err)
			}
			_, err = clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("error recreating Pod: %w", err)
			}
		} else {
			return fmt.Errorf("error creating Pod: %w", err)
		}
	}
	return nil
}

func createService(clientset *kubernetes.Clientset, targetName, namespace, uniqueID string) error {
	serviceName := fmt.Sprintf("%s-dast-report-%s-service", targetName, uniqueID)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": fmt.Sprintf("%s-dast-report-%s", targetName, uniqueID),
			},
			Ports: []corev1.ServicePort{
				{
					Port:     80,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
	_, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			// Update existing Service
			_, err = clientset.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("error updating existing Service: %w", err)
			}
		} else {
			return fmt.Errorf("error creating Service: %w", err)
		}
	}
	return nil
}

func sendWebhook(webhookURL, targetName, message, uniqueID string) {
	payload := map[string]string{
		"title":                  fmt.Sprintf("%s - DAST Scan", targetName),
		"description":            message,
		"external_aggregate_key": fmt.Sprintf("%s-dast-report-%s", targetName, uniqueID),
		"action":                 "alarmed",
		"severity":               "critical",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling webhook payload: %v\n", err)
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "1PcustomAuth/1.0",
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending webhook: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Webhook returned non-OK status: %s, response: %s\n", resp.Status, string(body))
	} else {
		fmt.Println("Webhook sent successfully.")
	}
}
