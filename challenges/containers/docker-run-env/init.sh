#!/bin/bash
docker pull mysql:latest 2>/dev/null
docker rm -f my-mysql 2>/dev/null
