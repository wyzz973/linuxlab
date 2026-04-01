#!/bin/bash
# Remove trailing empty line left by heredoc setup
sed -i '${/^$/d;}' challenge.txt
sed -i '${/^$/d;}' header.html
sed -i '/<!-- HEADER -->/r header.html' challenge.txt
sed -i '/<!-- HEADER -->/d' challenge.txt
