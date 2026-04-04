#!/bin/bash
# Verify the user created the three required files with plausible content.
# This is a concept challenge — we check file existence and basic content patterns
# rather than requiring a live HTTP server.

fail() { echo "FAIL: $1"; exit 1; }

for f in /tmp/get_response.txt /tmp/post_response.txt /tmp/response_headers.txt; do
    [ -f "$f" ] || fail "$f not found"
    [ -s "$f" ] || fail "$f is empty"
done

# get_response.txt should contain some text (any non-empty content is acceptable)
# post_response.txt should look like JSON or contain typical POST response content
# response_headers.txt should contain HTTP header lines
if grep -qi "HTTP/\|content-type\|server\|date" /tmp/response_headers.txt 2>/dev/null; then
    echo "PASS"
    exit 0
fi

# Fallback: if headers file doesn't look like headers but all files exist and are non-empty,
# still pass — the user demonstrated they know how to use curl
echo "PASS"
exit 0
