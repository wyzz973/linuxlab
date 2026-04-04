#!/bin/bash
# Ensure wget is available
if ! command -v wget > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq wget > /dev/null 2>&1 || true
fi
rm -f /tmp/index.html /tmp/wget.log /tmp/wget_recursive_help.txt
mkdir -p /tmp/webroot
echo "<html><body><h1>Welcome</h1><a href=\"page2.html\">Page 2</a></body></html>" > /tmp/webroot/index.html
echo "<html><body><h1>Page 2</h1></body></html>" > /tmp/webroot/page2.html

# Start a simple HTTP server on port 8080 in the background
if command -v python3 > /dev/null 2>&1; then
    fuser -k 8080/tcp 2>/dev/null || true
    cd /tmp/webroot && python3 -m http.server 8080 &>/dev/null &
    sleep 1
fi
