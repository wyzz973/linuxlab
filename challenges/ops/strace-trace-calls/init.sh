#!/bin/bash
apt-get update -qq && apt-get install -y -qq strace > /dev/null 2>&1
rm -f /tmp/strace_ls.txt /tmp/strace_open.txt /tmp/strace_summary.txt
