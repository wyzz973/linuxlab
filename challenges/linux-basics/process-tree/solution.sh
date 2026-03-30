#!/bin/bash
pstree -p > /tmp/result.txt 2>/dev/null || ps auxf > /tmp/result.txt
