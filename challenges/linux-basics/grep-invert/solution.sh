#!/bin/bash
grep -v 'grep' /home/lab/processes.txt | grep -v '\[' > /tmp/result.txt
