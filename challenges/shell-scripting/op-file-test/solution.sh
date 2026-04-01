#!/bin/bash
base="/home/learner/testdir"
[ -f "$base/regular.txt" ] && echo "is file" || echo "not file"
[ -d "$base/regular.txt" ] && echo "is directory" || echo "not directory"
[ -r "$base/regular.txt" ] && echo "is readable" || echo "not readable"
[ -s "$base/regular.txt" ] && echo "has content" || echo "is empty"
[ -d "$base/subdir" ] && echo "is directory" || echo "not directory"
[ -x "$base/script.sh" ] && echo "is executable" || echo "not executable"
[ -e "$base/nonexistent" ] && echo "exists" || echo "not exists"
