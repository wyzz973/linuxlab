#!/bin/bash
head -n 5 /home/lab/server.log > /tmp/head_result.txt
tail -n 3 /home/lab/server.log > /tmp/tail_result.txt
