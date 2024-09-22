package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("api-service"),
		)),
	)
	otel.SetTracerProvider(provider)

	return provider, nil
}

func main() {
	// Initialize OpenTelemetry tracer
	tracerProvider, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	// Create a new router
	r := mux.NewRouter()

	// Add a traced handler using middleware for all routes
	r.Use(otelMiddleware)

	// Register the API routes
	r.HandleFunc("/hello", helloHandler)

	// Wrap the router with OpenTelemetry instrumentation
	httpHandler := otelhttp.NewHandler(r, "server-handler")

	// Start the server
	log.Println("Starting server on :3333")
	log.Fatal(http.ListenAndServe(":3333", httpHandler))
}

// Middleware to automatically trace every request
func otelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("api-service")
		ctx, span := tracer.Start(r.Context(), r.URL.Path)
		defer span.End()

		// Pass the new context with the span to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// helloHandler is an example API handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the span from the context (if needed)
	// span := trace.SpanFromContext(r.Context())
	// span.AddEvent("Handling /hello request")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
