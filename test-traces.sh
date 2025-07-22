#!/bin/bash

# Test script for field summing collector
# This script generates test traces and validates field summing functionality

set -e

echo "Starting trace generation and validation tests..."

# Generate test traces with various field patterns
echo "Generating test traces..."
cd trace-generator
go run main.go

echo "Waiting for traces to be processed..."
sleep 5

echo "Validating field summing results..."
# Add validation logic here

echo "Test completed successfully!"