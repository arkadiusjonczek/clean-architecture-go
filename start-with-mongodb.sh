#!/bin/bash

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
DRIVER=mongodb go run ./cmd/server