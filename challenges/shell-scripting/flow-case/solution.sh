#!/bin/bash
while read file; do
    case "$file" in
        *.sh)
            echo "$file: 脚本文件"
            ;;
        *.txt|*.md)
            echo "$file: 文本文件"
            ;;
        *.jpg|*.png|*.gif)
            echo "$file: 图片文件"
            ;;
        *.tar.gz|*.zip)
            echo "$file: 压缩文件"
            ;;
        *)
            echo "$file: 未知类型"
            ;;
    esac
done < /home/learner/files.txt
