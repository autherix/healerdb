#!/usr/bin/env bash

# This script is used to run the docker container for the healerdb

# First check if docker is installed and running 
if ! docker info > /dev/null 2>&1; then
    echo "Docker is not installed or not running"
    exit 1
fi

# Check if the port 27017 is already in use
if netstat -tulpn | grep 27017 > /dev/null 2>&1; then
    echo "Port 27017 is already in use, make sure to stop the container or the service"
    exit 1
fi

# Check if the container is already running
if docker ps | grep healerdb > /dev/null 2>&1; then
    echo "Container healerdb is already running"
    exit 1
fi

docker run -d --name healerdb -p 27017:27017 mongo > /dev/null 2>&1

# Check if the container is running
if docker ps | grep healerdb > /dev/null 2>&1; then
    echo "Container healerdb is running"
    exit 0
else
    echo "Container healerdb is not running"
    exit 1
fi

echo "Exiting..."