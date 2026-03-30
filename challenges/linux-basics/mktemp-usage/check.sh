#!/bin/bash
if [ ! -f /tmp/tempfile_path.txt ] || [ ! -f /tmp/tempdir_path.txt ]; then
    echo "FAIL: Result files not found"
    exit 1
fi
tempfile=$(cat /tmp/tempfile_path.txt | tr -d '\n')
tempdir=$(cat /tmp/tempdir_path.txt | tr -d '\n')
if [ -f "$tempfile" ] && [ -d "$tempdir" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Temp file ($tempfile) or dir ($tempdir) not found"
    exit 1
fi
