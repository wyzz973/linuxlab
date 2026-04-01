#!/bin/bash
docker rmi tagdemo:v1 tagdemo:latest myregistry/tagdemo:v1 2>/dev/null
rm -rf /tmp/tagdemo
mkdir -p /tmp/tagdemo
rm -f /tmp/tag-list.txt
