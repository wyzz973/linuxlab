#!/bin/bash
docker pull alpine:latest 2>/dev/null
docker rm -f prune-test-1 prune-test-2 prune-test-3 2>/dev/null
# Create stopped containers
docker run --name prune-test-1 alpine echo "test1"
docker run --name prune-test-2 alpine echo "test2"
docker run --name prune-test-3 alpine echo "test3"
