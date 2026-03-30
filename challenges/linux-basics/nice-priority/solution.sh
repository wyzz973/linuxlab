#!/bin/bash
nice -n 10 sleep 600 &
echo $! > /tmp/result.txt
