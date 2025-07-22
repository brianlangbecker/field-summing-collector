#!/bin/bash

# Helper script to run the collector and test it

echo "Starting OTEL Collector in background..."
otelcol --config=trace-enrichment-config.yaml &
COLLECTOR_PID=$!

echo "Waiting 3 seconds for collector to start..."
sleep 3

echo "Sending test trace data..."
./test-traces.sh

echo ""
echo "Collector PID: $COLLECTOR_PID"
echo "To stop the collector, run: kill $COLLECTOR_PID"