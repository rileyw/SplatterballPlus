#!/bin/bash

# SplatServer Deployment Script

set -e

echo "🚀 Deploying SplatServer..."

# Build the application
echo "📦 Building application..."
go build -o splatserver .

# Run tests
echo "🧪 Running tests..."
DB_TYPE=sqlite go test -v ./...

# Build Docker image
echo "🐳 Building Docker image..."
docker build -t splatserver:latest .

echo "✅ Deployment preparation complete!"
echo ""
echo "To run locally:"
echo "  ./splatserver"
echo ""
echo "To run with Docker Compose:"
echo "  docker-compose up -d"
echo ""
echo "To deploy to cloud:"
echo "  1. Push image to container registry"
echo "  2. Deploy to Kubernetes/Docker Swarm"
echo "  3. Configure load balancer"
echo "  4. Set up monitoring"
