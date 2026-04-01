#!/bin/bash
docker network create frontend
docker network create backend
docker run -d --network frontend --name proxy-server nginx
docker run -d --network backend --name app-server alpine sleep 3600
docker network connect backend proxy-server
