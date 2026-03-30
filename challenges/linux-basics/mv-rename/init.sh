#!/bin/bash
mkdir -p /home/lab/temp /home/lab/archive
echo "This file should be renamed" > /home/lab/old_name.txt
echo "col1,col2,col3" > /home/lab/temp/data.csv
echo "1,2,3" >> /home/lab/temp/data.csv
