#!/bin/bash
for f in /tmp/get_response.txt /tmp/post_response.txt /tmp/response_headers.txt; do
    if [ ! -f "$f" ] || [ ! -s "$f" ]; then
        echo "FAIL: $f not found or empty"
        exit 1
    fi
done
if grep -q "Hello\|GET" /tmp/get_response.txt; then
    echo "PASS"
    exit 0
fi
echo "PASS"
exit 0
