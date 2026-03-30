#!/bin/bash
owner=$(stat -c %U /home/lab/webapp/index.html 2>/dev/null)
group=$(stat -c %G /home/lab/webapp/index.html 2>/dev/null)
if [ "$owner" = "www-data" ] && [ "$group" = "www-data" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected www-data:www-data, got $owner:$group"
    exit 1
fi
