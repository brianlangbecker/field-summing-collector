# Field Summing Collector Job

Complete automation setup for running the trace generator and OTEL collector pipeline.

## Quick Start

### Docker Compose (Recommended)
```bash
export HONEYCOMB_API_KEY="your-api-key"
./run-job.sh
```

### Manual Docker Build
```bash
docker build -t field-summing-collector .
docker run -e OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 field-summing-collector
```

## Job Options

### Environment Variables
- `HONEYCOMB_API_KEY`: Required for Honeycomb export
- `RUN_DURATION`: Seconds to run (default: continuous)
- `OTEL_EXPORTER_OTLP_ENDPOINT`: Collector endpoint

### Configuration (`job-config.yaml`)
- `traces_per_batch`: Number of traces to generate per batch
- `batch_interval`: Time between batches
- `total_batches`: Total number of batches (0 = infinite)

## Deployment Options

### 1. Local Development
```bash
./run-job.sh
```

### 2. Docker Compose Stack
```bash
docker-compose up -d
```

### 3. Kubernetes Job
```bash
kubectl apply -f k8s-job.yaml
```

### 4. Kubernetes CronJob (Every 5 minutes)
```bash
# Edit the secret with your API key first
kubectl apply -f k8s-job.yaml
```

## Services Included

- **OTEL Collector**: Processes and enriches traces
- **Trace Generator**: Generates realistic trace data
- **Prometheus**: Metrics collection
- **Grafana**: Visualization dashboard

## Monitoring

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Collector Metrics**: http://localhost:8888/metrics

## Scaling

Adjust in `job-config.yaml`:
```yaml
trace_generator:
  traces_per_batch: 20      # More traces per batch
  batch_interval: "10s"     # Faster generation
  total_batches: 100        # Finite run
```

## Troubleshooting

### Check logs:
```bash
docker-compose logs trace-generator
docker-compose logs otel-collector
```

### Test collector connectivity:
```bash
curl -X POST http://localhost:4318/v1/traces \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
```