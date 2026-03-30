#!/bin/bash
curl -sS -o /tmp/page.html http://localhost:8080/ 2>/dev/null || curl -sS http://localhost > /tmp/result.txt 2>&1 || touch /tmp/result.txt
