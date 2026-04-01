#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq >/dev/null 2>&1
apt-get install -y -qq --no-install-recommends wget python3-minimal >/dev/null 2>&1
mkdir -p /home/lab/www /home/lab/downloads
echo "<html><body><h1>Test Page</h1></body></html>" > /home/lab/www/index.html
# Start simple HTTP server
cd /home/lab/www && python3 -m http.server 8080 &>/dev/null &
sleep 1
rm -f /home/lab/downloads/index.html
