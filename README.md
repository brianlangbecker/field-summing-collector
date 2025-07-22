# OTEL Collector - Field Summing Collector

This project contains an OpenTelemetry collector configuration that sums field values from specific spans using the sum connector. It's designed to aggregate numeric fields (like duration) from filtered spans.

## Configuration

The collector configuration assumes a specific trace structure:
- **One `getData` operation per trace** from the `appserver` service
- **All `Signal/read` operations are descendants** of the `getData` span
- **Trace grouping** ensures proper processing boundaries

The collector configuration (`trace-enrichment-config.yaml`) includes:

- **groupbytrace processor**: Groups spans by trace ID for proper processing
- **transform processor**: Sums Signal/read durations and adds to getData spans
- **Batch processor**: Optimizes export performance

## Quick Start with Docker

### 1. Set up environment
```bash
# Set your Honeycomb API key
export HONEYCOMB_API_KEY="your-api-key-here"

# Clone or navigate to the project
cd field-summing-collector
```

### 2. Run the complete stack
```bash
# Start collector, trace generator, and monitoring
./run-job.sh
```

This will:
- Start the OTEL collector with trace enrichment
- Generate realistic trace data continuously
- Export enriched traces to Honeycomb
- Provide monitoring dashboards

### 3. Monitor the system
- **Grafana Dashboard**: http://localhost:3000 (admin/admin)
- **Prometheus Metrics**: http://localhost:9090
- **Collector Metrics**: http://localhost:8888/metrics

### 4. View results in Honeycomb
Look for traces with these added attributes on `getData` spans:
- `read_signal_sum_ms`: Sum of all Signal/read durations
- `read_signal_max_end_ns`: Maximum end time of Signal/read spans

## Manual Testing

### Option 1: Docker Compose
```bash
# Start services
docker-compose up -d

# Check logs
docker-compose logs trace-generator
docker-compose logs otel-collector

# Stop services
docker-compose down
```

### Option 2: Local Development
```bash
# Start collector
otelcol --config=trace-enrichment-config.yaml &

# In another terminal, run trace generator
cd trace-generator
go run main.go
```

### Option 3: Script-based testing
```bash
# Use the original curl script
./test-traces.sh
```

## Expected Output

The collector will:
- Process traces with `groupbytrace` processor
- Sum Signal/read durations from `datasource-proxy` service
- Add `read_signal_sum_ms` and `read_signal_max_end_ns` attributes to `getData` spans
- Export enriched traces to Honeycomb

## Troubleshooting

### Check collector connectivity:
```bash
curl -X POST http://localhost:4318/v1/traces \
  -H "Content-Type: application/json" \
  -d '{"test": "connectivity"}'
```

### View collector logs:
```bash
docker-compose logs otel-collector
```

### Verify trace generation:
```bash
docker-compose logs trace-generator
```

## Customization

To modify for different services or fields, edit `trace-enrichment-config.yaml`:
- Change service names in the Select conditions
- Modify span name filters
- Adjust attribute names for the enriched fields