#!/bin/bash
cut -d':' -f1,7 /home/lab/passwd_sample.txt | tr 'A-Z' 'a-z' > /tmp/result.txt
