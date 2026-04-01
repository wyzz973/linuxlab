#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f lifecycle-test 2>/dev/null
