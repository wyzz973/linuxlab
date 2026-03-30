#!/bin/bash
du -s /home/lab/*/ | sort -rn | head -3 > /tmp/result.txt
