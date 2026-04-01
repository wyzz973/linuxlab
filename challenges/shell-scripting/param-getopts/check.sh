#!/bin/bash
out1=$(bash /home/learner/solution.sh -n Alice -a 25 -v 2>/dev/null)
expected1=$(printf "Verbose mode enabled\nName: Alice\nAge: 25")
out2=$(bash /home/learner/solution.sh -n Bob -a 30 2>/dev/null)
expected2=$(printf "Name: Bob\nAge: 30")
if [ "$out1" = "$expected1" ] && [ "$out2" = "$expected2" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "Test 1: '$out1'"
    echo "Test 2: '$out2'"
    exit 1
fi
