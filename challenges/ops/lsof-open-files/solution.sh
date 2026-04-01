#!/bin/bash
lsof -i TCP -sTCP:LISTEN > /tmp/lsof_listen.txt 2>&1
lsof -p 1 > /tmp/lsof_pid.txt 2>&1
lsof 2>/dev/null | grep "(deleted)" > /tmp/lsof_deleted.txt || echo "No deleted files" > /tmp/lsof_deleted.txt
