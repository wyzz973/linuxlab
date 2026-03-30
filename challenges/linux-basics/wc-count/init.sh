#!/bin/bash
mkdir -p /home/lab
python3 -c "
for i in range(1, 21):
    print(f'Line {i}: This is sample text for the article about Linux commands.')
" > /home/lab/article.txt 2>/dev/null || {
    for i in $(seq 1 20); do
        echo "Line $i: This is sample text for the article about Linux commands."
    done > /home/lab/article.txt
}
