#!/bin/bash
for f in /tmp/loaded_services.txt /tmp/failed_units.txt /tmp/unit_files.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
if [ -s /tmp/loaded_services.txt ] || [ -s /tmp/unit_files.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: output files are empty"
exit 1
