#!/bin/bash
mkdir -p /home/lab/webapp/{css,js,img}
echo "body { margin: 0; }" > /home/lab/webapp/css/style.css
echo "console.log('hello');" > /home/lab/webapp/js/app.js
echo "<html><body>Hello</body></html>" > /home/lab/webapp/index.html
echo "PNG_PLACEHOLDER" > /home/lab/webapp/img/logo.png
rm -rf /home/lab/webapp-backup
