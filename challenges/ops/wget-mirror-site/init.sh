#!/bin/bash
apt-get update -qq && apt-get install -y -qq wget python3 > /dev/null 2>&1
rm -f /tmp/index.html /tmp/wget.log /tmp/wget_recursive_help.txt
mkdir -p /tmp/webroot
echo "<html><body><h1>Welcome</h1><a href=\"page2.html\">Page 2</a></body></html>" > /tmp/webroot/index.html
echo "<html><body><h1>Page 2</h1></body></html>" > /tmp/webroot/page2.html
cd /tmp/webroot && python3 -m http.server 8080 &>/dev/null &
sleep 1
