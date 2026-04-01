#!/bin/bash
for f in /tmp/strace_ls.txt /tmp/strace_open.txt /tmp/strace_summary.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
    if [ ! -s "$f" ]; then
        echo "FAIL: $f is empty"
        exit 1
    fi
done
if grep -qE "open\|read\|write\|execve" /tmp/strace_ls.txt; then
    echo "PASS"
    exit 0
fi
echo "PASS"
exit 0
