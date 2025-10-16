#!/bin/bash

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "Stopping MongoDB service..."
    docker-compose down
    echo "Cleanup complete."
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM EXIT

# Start MongoDB service
echo "Starting MongoDB service..."
docker-compose up -d mongodb

# Wait for MongoDB to be ready
echo "Waiting for MongoDB to be ready..."
until docker-compose exec mongodb mongosh --eval "db.runCommand('ping')" > /dev/null 2>&1; do
    echo "Waiting for MongoDB..."
    sleep 2
done

echo "MongoDB is ready!"

# Start the application with MongoDB driver
echo "Starting application with MongoDB driver..."
echo "Press Ctrl+C to stop the application and MongoDB service"
DRIVER=mongodb go run ./cmd/server

# If we reach here, the application exited normally
echo "Application stopped normally."