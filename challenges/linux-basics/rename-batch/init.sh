#!/bin/bash
mkdir -p /home/lab/images
for i in 001 002 003 004 005; do
    echo "JPEG_DATA_$i" > "/home/lab/images/IMG_$i.JPG"
done
