#!/bin/bash
groupadd developers 2>/dev/null || true
mkdir -p /home/lab/project/src
echo "print('hello')" > /home/lab/project/src/main.py
echo "# README" > /home/lab/project/README.md
