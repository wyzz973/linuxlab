#!/bin/bash
mkdir -p /home/lab/webroot/{css,js,images}
echo "<html></html>" > /home/lab/webroot/index.html
echo "body{}" > /home/lab/webroot/css/style.css
echo "var x;" > /home/lab/webroot/js/app.js
chmod -R 600 /home/lab/webroot
