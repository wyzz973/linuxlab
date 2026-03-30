#!/bin/bash
useradd -r www-data 2>/dev/null || true
mkdir -p /home/lab/webapp
echo "<html>Hello</html>" > /home/lab/webapp/index.html
echo "config=true" > /home/lab/webapp/config.ini
