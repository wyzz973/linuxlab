#!/bin/bash
docker rmi multistage:v1 2>/dev/null
rm -rf /tmp/multistage
mkdir -p /tmp/multistage
