#!/bin/bash

# Variables
PROJECT_DIR=~/timewise-dms
DOCKER_IMAGE_NAME="timewise-dms"
DOCKER_CONTAINER_NAME="timewise-dms-container"
PORT_MAPPING="3030:3030"

# Step 1: Navigate to the project directory
echo "Navigating to the project directory..."
cd "$PROJECT_DIR" || { echo "Directory $PROJECT_DIR not found!"; exit 1; }

# Step 2: Pull the latest code from GitHub
echo "Pulling the latest code from GitHub..."
git pull origin main || { echo "Git pull failed!"; exit 1; }

# Step 3: Build the Docker image
echo "Building the Docker image..."
sudo docker build -t "$DOCKER_IMAGE_NAME" . || { echo "Docker build failed!"; exit 1; }

# Step 4: Stop and remove the existing container
echo "Stopping the existing container (if running)..."
if sudo docker ps -q --filter "name=$DOCKER_CONTAINER_NAME" | grep -q .; then
    sudo docker stop "$DOCKER_CONTAINER_NAME" || { echo "Failed to stop the container!"; exit 1; }
    sudo docker rm "$DOCKER_CONTAINER_NAME" || { echo "Failed to remove the container!"; exit 1; }
else
    echo "No running container found with the name $DOCKER_CONTAINER_NAME."
fi

# Step 5: Run the new container
echo "Starting the new container..."
sudo docker run -d --name "$DOCKER_CONTAINER_NAME" -p "$PORT_MAPPING" "$DOCKER_IMAGE_NAME" || { echo "Failed to start the container!"; exit 1; }

# Step 6: Verify the container is running
echo "Verifying the container is running..."
sudo docker ps | grep "$DOCKER_CONTAINER_NAME" && echo "Deployment successful!" || echo "Deployment failed!"
