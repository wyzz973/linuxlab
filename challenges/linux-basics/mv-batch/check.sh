#!/bin/bash
for f in vacation.jpg selfie.jpg family.jpg; do
    if [ ! -f "/home/lab/photos/$f" ]; then
        echo "FAIL: $f not found in /home/lab/photos/"
        exit 1
    fi
done
# Non-jpg files should stay
if [ ! -f /home/lab/inbox/document.pdf ]; then
    echo "FAIL: document.pdf should still be in inbox"
    exit 1
fi
echo "PASS"
exit 0
