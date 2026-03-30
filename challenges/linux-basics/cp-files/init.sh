#!/bin/bash
mkdir -p /home/lab/source /home/lab/backup
echo "Monthly Report - Q4 2025" > /home/lab/source/report.txt
echo "Revenue: $1,200,000" >> /home/lab/source/report.txt
echo "Expenses: $800,000" >> /home/lab/source/report.txt
chmod 644 /home/lab/source/report.txt
touch -t 202501150900 /home/lab/source/report.txt
