#!/bin/bash

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
DRIVER=mysql go run ./cmd/server