#!/bin/bash

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "Stopping MySQL service..."
    docker-compose down
    echo "Cleanup complete."
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM EXIT

# Start MySQL service
echo "Starting MySQL service..."
docker-compose up -d mysql

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
until docker-compose exec mysql mysqladmin ping -h localhost --silent; do
    echo "Waiting for MySQL..."
    sleep 2
done

echo "MySQL is ready!"

# Start the application with MySQL driver
echo "Starting application with MySQL driver..."
echo "Press Ctrl+C to stop the application and MySQL service"
DRIVER=mysql go run ./cmd/server

# If we reach here, the application exited normally
echo "Application stopped normally."