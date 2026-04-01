#!/bin/bash
docker pull nginx:latest alpine:latest 2>/dev/null
docker rm -f running-web stopped-task 2>/dev/null
rm -f /tmp/all-containers.txt /tmp/running-containers.txt
