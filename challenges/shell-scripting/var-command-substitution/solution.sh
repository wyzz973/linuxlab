#!/bin/bash
today=$(date +%Y-%m-%d)
user=$(whoami)
lines=$(wc -l < /etc/passwd)
echo "Today: $today"
echo "User: $user"
echo "Lines in passwd: $lines"
