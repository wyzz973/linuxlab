#!/bin/bash
wc -l < /home/learner/numbers.txt | tr -d ' '
sort -n /home/learner/numbers.txt
sort -n /home/learner/numbers.txt | tail -1
sort -n /home/learner/numbers.txt | head -1
