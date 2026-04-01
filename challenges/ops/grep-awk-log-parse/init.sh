#!/bin/bash
rm -f /tmp/errors.txt /tmp/error_summary.txt /tmp/hourly_errors.txt
cat > /var/log/app_sim.log << LOGEOF
2024-03-15 08:15:23 INFO  [main] Application started
2024-03-15 08:15:25 DEBUG [db] Connection pool initialized
2024-03-15 08:30:10 ERROR [api] Failed to connect to database: timeout
2024-03-15 09:00:05 INFO  [api] Request processed successfully
2024-03-15 09:15:30 ERROR [auth] Invalid token for user admin
2024-03-15 09:45:20 WARN  [api] Slow query detected: 3.5s
2024-03-15 10:00:15 ERROR [api] NullPointerException in UserService
2024-03-15 10:30:45 ERROR [db] Connection pool exhausted
2024-03-15 10:45:00 INFO  [api] Health check OK
2024-03-15 11:00:10 ERROR [api] Request timeout after 30s
2024-03-15 11:30:20 WARN  [mem] Memory usage above 80%
2024-03-15 11:45:30 ERROR [api] Failed to serialize response
2024-03-15 12:00:05 INFO  [gc] Garbage collection completed
2024-03-15 08:45:15 ERROR [net] Connection refused to upstream
LOGEOF
