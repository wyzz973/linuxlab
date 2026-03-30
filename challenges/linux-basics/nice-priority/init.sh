#!/bin/bash
pkill -f "sleep 600" 2>/dev/null || true
rm -f /tmp/result.txt
