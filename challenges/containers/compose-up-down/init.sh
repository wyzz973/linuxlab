#!/bin/bash
docker pull nginx:alpine redis:7-alpine 2>/dev/null
docker rm -f updown-web updown-cache 2>/dev/null
rm -f /tmp/compose-status.txt
