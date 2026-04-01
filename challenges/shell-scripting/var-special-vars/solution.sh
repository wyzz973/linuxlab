#!/bin/bash
echo "Script name: $0"
echo "PID: $$"
ls /nonexistent 2>/dev/null
echo "Exit status: $?"
echo "test" > /dev/null
echo "Exit status: $?"
