#!/bin/bash
rm -f /tmp/page.html /tmp/result.txt
# Start a simple web server for the challenge
mkdir -p /home/lab/www
echo "<html><body><h1>Welcome to LinuxLab</h1></body></html>" > /home/lab/www/index.html
cd /home/lab/www && python3 -m http.server 8080 &>/dev/null &
sleep 1
