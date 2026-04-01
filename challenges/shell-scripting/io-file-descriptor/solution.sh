#!/bin/bash
exec 3>/home/learner/fd_output.txt
echo "Written via fd 3" >&3
exec 3>&-
cat /home/learner/fd_output.txt

exec 0</home/learner/fd_input.txt
read line
echo "$line"
