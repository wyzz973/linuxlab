#!/bin/bash
# Remove trailing empty line left by heredoc setup
sed -i '${/^$/d;}' challenge.txt
sed -i "s/'/\"/g" challenge.txt
