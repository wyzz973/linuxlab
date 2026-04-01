#!/bin/bash
while read line; do
    echo "$line: ${#line}"
done < /home/learner/words.txt
