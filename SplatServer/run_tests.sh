#!/bin/bash

# run_tests.sh: Script to run all tests for SplatServer

echo "Running unit tests..."
go test -v -timeout 30s
if [ $? -ne 0 ]; then
    echo "Unit tests failed!"
    exit 1
fi

echo "Starting server in background..."
# Kill any existing process on port 4000
lsof -ti:4000 | xargs kill -9 2>/dev/null || true
DB_TYPE=sqlite go run . &
SERVER_PID=$!
sleep 2  # Wait for server to start

echo "Running integration test client..."
cd client && go run test_client.go
if [ $? -ne 0 ]; then
    echo "Integration test failed!"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

echo "Stopping server..."
kill $SERVER_PID 2>/dev/null

echo "All tests passed!"
