#!/bin/bash
grep -nE '[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+' /home/lab/emails.txt > /tmp/result.txt
