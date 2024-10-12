#!/bin/bash

# Set variables
IMAGE_NAME="website"
CONTAINER_NAME="forum"

# Build the Docker image
echo "Building Docker image..."
docker build -t $IMAGE_NAME .

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Docker image built successfully."
    
    # Stop and remove any existing container with the same name
    docker stop $CONTAINER_NAME 2>/dev/null
    docker rm $CONTAINER_NAME 2>/dev/null
    
    # Run the container with explicit port mapping
    echo "Running Docker container..."
    docker run -d --name $CONTAINER_NAME -p 8989:8989 $IMAGE_NAME
    
    echo "Container '$CONTAINER_NAME' is now running. You can access it at http://localhost:8080"
    
    # Verify the container is running
    docker ps | grep $CONTAINER_NAME
else
    echo "Failed to build Docker image."
    exit 1
fi