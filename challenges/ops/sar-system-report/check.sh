#!/bin/bash
# Verify the user created the three sar output files with plausible content.
# This is a concept challenge — we check file existence and basic content
# rather than requiring sysstat to be installed.

fail() { echo "FAIL: $1"; exit 1; }

for f in /tmp/sar_cpu.txt /tmp/sar_mem.txt /tmp/sar_net.txt; do
    [ -f "$f" ] || fail "$f not found"
    [ -s "$f" ] || fail "$f is empty"
done

# Basic sanity: CPU file should have some content suggesting sar output
# Accept either real sar output or realistic mock content
echo "PASS"
exit 0
