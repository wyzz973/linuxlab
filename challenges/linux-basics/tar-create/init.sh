#!/bin/bash
mkdir -p /home/lab/website/{css,js,images}
echo "<html><body><h1>Welcome</h1></body></html>" > /home/lab/website/index.html
echo "body { font-family: Arial; }" > /home/lab/website/css/style.css
echo "alert('hello');" > /home/lab/website/js/main.js
echo "PNG_DATA" > /home/lab/website/images/banner.png
rm -f /home/lab/website-backup.tar.gz
