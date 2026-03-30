#!/bin/bash
mkdir -p /home/lab/temp_mystery
echo "You found the secret!" > /home/lab/temp_mystery/secret.txt
echo "Not so secret" > /home/lab/temp_mystery/readme.txt
echo "data data data" > /home/lab/temp_mystery/data.bin
tar czf /home/lab/mystery.tar.gz -C /home/lab/temp_mystery .
rm -rf /home/lab/temp_mystery
