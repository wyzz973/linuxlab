#!/bin/bash
mount | grep -E 'ext4|xfs' > /tmp/result.txt 2>/dev/null
if [ ! -s /tmp/result.txt ]; then
    mount > /tmp/result.txt
fi
