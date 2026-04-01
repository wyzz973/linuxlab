#!/bin/bash
strace ls / > /dev/null 2> /tmp/strace_ls.txt
grep -E "^open|^openat" /tmp/strace_ls.txt > /tmp/strace_open.txt
strace -c ls / > /dev/null 2> /tmp/strace_summary.txt
