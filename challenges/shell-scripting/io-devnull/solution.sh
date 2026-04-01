#!/bin/bash
ls /nonexistent 2>/dev/null
echo "Error suppressed"
ls /home >/dev/null
echo "Output suppressed"
ls /nonexistent 2>/home/learner/error.log
wc -l < /home/learner/error.log | tr -d ' '
ls /nonexistent >/dev/null 2>&1
echo "All done"
