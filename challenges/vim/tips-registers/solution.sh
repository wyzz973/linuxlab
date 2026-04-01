#!/bin/bash
# Remove trailing empty line left by heredoc setup
sed -i '${/^$/d;}' challenge.txt
printf "gamma [1]\nalpha [2]\nbeta [3]\n" >> challenge.txt
