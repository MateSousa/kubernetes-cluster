FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o report-handler .

FROM alpine:3.20.3

COPY --from=builder /app/report-handler /usr/local/bin/

ENTRYPOINT ["report-handler"]

