#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/syslog.txt << 'EOF'
2025-01-01 10:00:01 INFO: System started
2025-01-01 10:00:02 INFO: Loading modules
2025-01-01 10:01:00 WARNING: High memory usage
2025-01-01 10:02:00 ERROR: Database connection failed
2025-01-01 10:02:01 INFO: Retrying connection
2025-01-01 10:02:05 ERROR: Connection timeout after 5s
2025-01-01 10:03:00 INFO: Connection restored
2025-01-01 10:04:00 WARNING: Disk space low
2025-01-01 10:05:00 ERROR: Failed to write log file
2025-01-01 10:06:00 INFO: Cleanup completed
EOF
