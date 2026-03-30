#!/bin/bash
dpkg -l > /tmp/installed.txt
apt search nginx > /tmp/search.txt 2>/dev/null || echo "apt search not available" > /tmp/search.txt
