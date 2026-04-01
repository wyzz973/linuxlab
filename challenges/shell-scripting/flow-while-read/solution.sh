#!/bin/bash
tail -n +2 /home/learner/data.csv | while IFS=',' read name age city; do
    echo "$name is $age years old, lives in $city"
done
