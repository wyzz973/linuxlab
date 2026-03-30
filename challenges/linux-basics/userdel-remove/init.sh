#!/bin/bash
userdel -r exemployee 2>/dev/null || true
useradd -m exemployee
echo "personal data" > /home/exemployee/data.txt
