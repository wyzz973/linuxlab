#!/bin/bash
docker pull mysql:8.0 2>/dev/null
docker rm -f compose-mysql 2>/dev/null
rm -rf /tmp/compose-env
mkdir -p /tmp/compose-env
