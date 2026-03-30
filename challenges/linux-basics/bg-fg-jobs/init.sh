#!/bin/bash
# Kill any existing sleep 300
pkill -f "sleep 300" 2>/dev/null || true
rm -f /tmp/result.txt
