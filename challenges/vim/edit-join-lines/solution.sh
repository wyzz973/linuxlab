#!/bin/bash
# Remove trailing empty line left by heredoc setup
sed -i '${/^$/d;}' challenge.txt
paste -sd ' ' challenge.txt > /tmp/vim_join.txt && mv /tmp/vim_join.txt challenge.txt
