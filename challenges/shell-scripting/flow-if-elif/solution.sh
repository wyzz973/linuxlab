#!/bin/bash
while read score; do
    if [ "$score" -ge 90 ]; then
        grade="A"
    elif [ "$score" -ge 80 ]; then
        grade="B"
    elif [ "$score" -ge 70 ]; then
        grade="C"
    elif [ "$score" -ge 60 ]; then
        grade="D"
    else
        grade="F"
    fi
    echo "$score: $grade"
done < /home/learner/scores.txt
