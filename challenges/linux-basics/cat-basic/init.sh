#!/bin/bash
mkdir -p /home/lab
echo "Part 1 content" > /home/lab/part1.txt
echo "First section data" >> /home/lab/part1.txt
echo "Part 2 content" > /home/lab/part2.txt
echo "Second section data" >> /home/lab/part2.txt
rm -f /home/lab/combined.txt
