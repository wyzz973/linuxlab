#!/bin/bash
docker pull redis:7-alpine
docker pull nginx:1.25
docker images > /tmp/image-list.txt
