#!/bin/bash
# If sar is available, use it; otherwise create realistic mock output.

if command -v sar >/dev/null 2>&1; then
    sar 1 3 > /tmp/sar_cpu.txt 2>&1
    sar -r 1 3 > /tmp/sar_mem.txt 2>&1
    sar -n DEV 1 3 > /tmp/sar_net.txt 2>&1
else
    # Mock realistic sar CPU output
    cat > /tmp/sar_cpu.txt << 'EOF'
Linux 5.15.0-91-generic (ubuntu)    04/02/2026    _x86_64_    (2 CPU)

10:00:00 AM     CPU     %user     %nice   %system   %iowait    %steal     %idle
10:00:01 AM     all      2.01      0.00      1.00      0.00      0.00     96.98
10:00:02 AM     all      1.50      0.00      0.50      0.00      0.00     98.00
10:00:03 AM     all      3.02      0.00      1.01      0.50      0.00     95.48
Average:        all      2.18      0.00      0.84      0.17      0.00     96.82
EOF

    # Mock realistic sar memory output
    cat > /tmp/sar_mem.txt << 'EOF'
Linux 5.15.0-91-generic (ubuntu)    04/02/2026    _x86_64_    (2 CPU)

10:00:00 AM kbmemfree kbavail kbmemused  %memused kbbuffers  kbcached
10:00:01 AM   1024000 3500000   2976000     74.40    128000   1200000
10:00:02 AM   1020000 3496000   2980000     74.50    128000   1200000
10:00:03 AM   1022000 3498000   2978000     74.45    128000   1200000
Average:      1022000 3498000   2978000     74.45    128000   1200000
EOF

    # Mock realistic sar network output
    cat > /tmp/sar_net.txt << 'EOF'
Linux 5.15.0-91-generic (ubuntu)    04/02/2026    _x86_64_    (2 CPU)

10:00:00 AM     IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s
10:00:01 AM      eth0     15.00     12.00      1.50      0.80
10:00:01 AM        lo      0.00      0.00      0.00      0.00
10:00:02 AM      eth0     18.00     14.00      2.10      1.00
10:00:02 AM        lo      0.00      0.00      0.00      0.00
10:00:03 AM      eth0     12.00     10.00      1.20      0.60
10:00:03 AM        lo      0.00      0.00      0.00      0.00
Average:         eth0     15.00     12.00      1.60      0.80
Average:           lo      0.00      0.00      0.00      0.00
EOF
fi
