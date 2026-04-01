#!/bin/bash
docker volume create app-data
docker volume create db-data
docker volume inspect app-data > /tmp/volume-info.txt
docker volume ls > /tmp/volume-list.txt
