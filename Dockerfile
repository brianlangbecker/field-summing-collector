FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY trace-generator/go.mod trace-generator/go.sum ./
RUN go mod download

COPY trace-generator/main.go ./
RUN go build -o trace-generator main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/trace-generator .
COPY job-config.yaml .

CMD ["./trace-generator"]