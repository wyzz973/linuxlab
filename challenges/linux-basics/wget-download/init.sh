#!/bin/bash
mkdir -p /home/lab/www /home/lab/downloads
echo "<html><body><h1>Test Page</h1></body></html>" > /home/lab/www/index.html
# Start simple HTTP server
cd /home/lab/www && python3 -m http.server 8080 &>/dev/null &
sleep 1
rm -f /home/lab/downloads/index.html
