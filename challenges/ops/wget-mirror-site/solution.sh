#!/bin/bash
wget -O /tmp/index.html -o /tmp/wget.log http://localhost:8080/
wget --help 2>&1 | grep -iA5 "recurs" > /tmp/wget_recursive_help.txt
