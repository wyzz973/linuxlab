#!/bin/bash
curl -s -H "User-Agent: LinuxLabAgent/1.0" http://localhost:8080/ > /tmp/custom_ua.txt
curl -s -c /tmp/cookies.txt http://localhost:8080/ > /dev/null
curl -sv http://localhost:8080/ > /tmp/verbose.txt 2>&1
