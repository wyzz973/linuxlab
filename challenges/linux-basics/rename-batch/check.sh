#!/bin/bash
# Check that .jpg files exist
count_jpg=$(ls /home/lab/images/*.jpg 2>/dev/null | wc -l)
count_JPG=$(ls /home/lab/images/*.JPG 2>/dev/null | wc -l)
if [ "$count_jpg" -eq 5 ] && [ "$count_JPG" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 5 .jpg files and 0 .JPG files (got $count_jpg .jpg and $count_JPG .JPG)"
    exit 1
fi
