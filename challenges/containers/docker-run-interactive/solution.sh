#!/bin/bash
docker run --name interactive-ubuntu ubuntu bash -c 'mkdir -p /data && echo "interactive mode" > /data/hello.txt'
