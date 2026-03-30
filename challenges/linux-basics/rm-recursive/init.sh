#!/bin/bash
mkdir -p /home/lab/old_project/{src,docs,build}
echo "old code" > /home/lab/old_project/src/main.c
echo "old docs" > /home/lab/old_project/docs/readme.txt
mkdir -p /home/lab/current_project
echo "important code" > /home/lab/current_project/app.py
