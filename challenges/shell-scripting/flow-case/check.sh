#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "setup.sh: 脚本文件\nreadme.md: 文本文件\nphoto.jpg: 图片文件\narchive.tar.gz: 压缩文件\ndata.csv: 未知类型")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
