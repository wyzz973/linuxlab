#!/bin/bash
while read name age role; do
    if [ "$age" -ge 18 ]; then
        if [ "$role" = "admin" ]; then
            echo "$name: full access"
        else
            echo "$name: limited access"
        fi
    else
        echo "$name: no access (underage)"
    fi
done < /home/learner/users.txt
