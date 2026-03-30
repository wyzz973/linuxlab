#!/bin/bash
mkdir -p /home/lab/downloads
wget -P /home/lab/downloads/ http://localhost:8080/index.html 2>/dev/null || echo "wget -P /home/lab/downloads/ http://localhost:8080/index.html" > /tmp/result.txt
