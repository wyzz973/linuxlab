#!/bin/bash
mkdir -p /home/lab/temp_data/data
echo "id,name,email" > /home/lab/temp_data/data/users.csv
echo "1,Alice,alice@example.com" >> /home/lab/temp_data/data/users.csv
echo "2,Bob,bob@example.com" >> /home/lab/temp_data/data/users.csv
echo "config=true" > /home/lab/temp_data/data/settings.ini
tar czf /home/lab/archive.tar.gz -C /home/lab/temp_data data
rm -rf /home/lab/temp_data /home/lab/extracted
