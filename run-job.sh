#!/bin/bash

# Field Summing Collector Job Runner
# This script runs the complete trace generation and processing pipeline

set -e

echo "🚀 Starting Field Summing Collector Job..."

# Check if HONEYCOMB_API_KEY is set
if [ -z "$HONEYCOMB_API_KEY" ]; then
    echo "❌ Error: HONEYCOMB_API_KEY environment variable is required"
    echo "Set it with: export HONEYCOMB_API_KEY='your-api-key'"
    exit 1
fi

# Function to cleanup on exit
cleanup() {
    echo "🧹 Cleaning up..."
    docker-compose down
}
trap cleanup EXIT

# Build and start services
echo "🔧 Building and starting services..."
docker-compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check service health
echo "🔍 Checking service health..."
docker-compose ps

# Show logs
echo "📋 Service logs:"
docker-compose logs --tail=50

# Keep running and show trace generator logs
echo "🔄 Monitoring trace generation..."
docker-compose logs -f trace-generator

# Run for specific duration if provided
if [ ! -z "$RUN_DURATION" ]; then
    echo "⏰ Running for $RUN_DURATION seconds..."
    sleep $RUN_DURATION
    echo "✅ Job completed after $RUN_DURATION seconds"
else
    echo "🔄 Running continuously. Press Ctrl+C to stop."
    # Keep running until interrupted
    while true; do
        sleep 30
        echo "📊 Status check..."
        docker-compose ps
    done
fi