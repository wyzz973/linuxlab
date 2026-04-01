#!/bin/bash
rm -f /tmp/tail_output.txt /tmp/live_errors.txt /tmp/line_count.txt
# Create a simulated live log
cat > /var/log/live_sim.log << LOGEOF
2024-03-15 10:00:01 INFO Starting service
2024-03-15 10:00:02 INFO Loading configuration
2024-03-15 10:00:03 DEBUG Config loaded successfully
2024-03-15 10:00:04 INFO Listening on port 8080
2024-03-15 10:00:05 ERROR Failed to bind to port 9090
2024-03-15 10:00:06 WARN Memory usage high
2024-03-15 10:00:07 INFO Processing request 1
2024-03-15 10:00:08 INFO Processing request 2
2024-03-15 10:00:09 ERROR Database connection lost
2024-03-15 10:00:10 INFO Reconnecting to database
2024-03-15 10:00:11 INFO Connection restored
2024-03-15 10:00:12 DEBUG Query executed in 0.5s
2024-03-15 10:00:13 INFO Processing request 3
2024-03-15 10:00:14 WARN Slow response: 2.1s
2024-03-15 10:00:15 INFO Processing request 4
2024-03-15 10:00:16 ERROR Timeout waiting for response
2024-03-15 10:00:17 INFO Retry successful
2024-03-15 10:00:18 INFO Processing request 5
2024-03-15 10:00:19 DEBUG Cache hit ratio: 85%
2024-03-15 10:00:20 INFO Health check passed
2024-03-15 10:00:21 INFO Processing request 6
2024-03-15 10:00:22 ERROR Invalid input from client
2024-03-15 10:00:23 INFO Request rejected
2024-03-15 10:00:24 INFO Processing request 7
2024-03-15 10:00:25 INFO Shutdown signal received
LOGEOF
# Background process to append lines
(sleep 1; echo "2024-03-15 10:01:00 ERROR New error appeared" >> /var/log/live_sim.log) &
