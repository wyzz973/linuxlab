#!/bin/bash
docker rmi greeter:v1 2>/dev/null
rm -rf /tmp/greeter
mkdir -p /tmp/greeter
