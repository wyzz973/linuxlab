#!/bin/bash
# Clean up any previous attempt
rm -f /tmp/sar_cpu.txt /tmp/sar_mem.txt /tmp/sar_net.txt

# Create a hint file for the user
cat > /tmp/README_challenge.txt << 'EOF'
Challenge: sar System Activity Report

The sar tool may not be installed in this environment.
Demonstrate your knowledge by creating the expected output files
with realistic sar-style content.

1. /tmp/sar_cpu.txt - CPU usage report (as from: sar 1 3)
2. /tmp/sar_mem.txt - Memory report (as from: sar -r 1 3)
3. /tmp/sar_net.txt - Network stats (as from: sar -n DEV 1 3)
EOF
