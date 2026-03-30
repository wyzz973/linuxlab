#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/large_log.txt << 'EOF'
2025-01-01 10:00:00 INFO Starting application
2025-01-01 10:00:01 INFO Loading modules
2025-01-01 10:00:02 INFO Database connected
2025-01-01 10:00:03 WARNING High memory usage detected
2025-01-01 10:00:04 INFO Attempting recovery
2025-01-01 10:00:05 CRITICAL Out of memory - service degraded
2025-01-01 10:00:06 ERROR Failed to allocate memory
2025-01-01 10:00:07 INFO Emergency cleanup initiated
2025-01-01 10:00:08 INFO Memory freed
2025-01-01 10:00:09 INFO Normal operation resumed
2025-01-01 10:01:00 INFO Processing requests
2025-01-01 10:01:01 INFO Request handled successfully
2025-01-01 10:02:00 WARNING Disk space low
2025-01-01 10:02:01 CRITICAL Disk full - writes failing
2025-01-01 10:02:02 ERROR Write operation failed
2025-01-01 10:02:03 INFO Disk cleanup started
EOF
