#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/data.txt << 'EOF'
Test 1: PASS - Basic functionality
Test 2: FAIL - Memory leak detected
Test 3: PASS - API response correct
Test 4: FAIL - Timeout exceeded
Test 5: PASS - Data integrity check
Test 6: PASS - Security scan clean
Test 7: FAIL - Permission denied
EOF
rm -f /tmp/passed.txt
