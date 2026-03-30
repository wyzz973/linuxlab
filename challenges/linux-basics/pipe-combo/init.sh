#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/web_access.log << 'EOF'
192.168.1.100 - - [01/Jan/2025:10:00:01] "GET /index.html HTTP/1.1" 200 1234
10.0.0.50 - - [01/Jan/2025:10:00:02] "GET /api/users HTTP/1.1" 200 567
192.168.1.100 - - [01/Jan/2025:10:00:03] "POST /api/login HTTP/1.1" 200 89
172.16.0.1 - - [01/Jan/2025:10:00:04] "GET /index.html HTTP/1.1" 200 1234
10.0.0.50 - - [01/Jan/2025:10:00:05] "GET /api/products HTTP/1.1" 200 2345
192.168.1.100 - - [01/Jan/2025:10:00:06] "GET /css/style.css HTTP/1.1" 200 456
192.168.1.200 - - [01/Jan/2025:10:00:07] "GET /index.html HTTP/1.1" 200 1234
10.0.0.50 - - [01/Jan/2025:10:00:08] "GET /api/orders HTTP/1.1" 200 789
192.168.1.100 - - [01/Jan/2025:10:00:09] "GET /images/logo.png HTTP/1.1" 200 5678
172.16.0.1 - - [01/Jan/2025:10:00:10] "POST /api/data HTTP/1.1" 201 123
192.168.1.100 - - [01/Jan/2025:10:01:01] "GET /api/status HTTP/1.1" 200 45
10.0.0.50 - - [01/Jan/2025:10:01:02] "GET /index.html HTTP/1.1" 200 1234
10.0.0.1 - - [01/Jan/2025:10:01:03] "GET /index.html HTTP/1.1" 200 1234
192.168.1.100 - - [01/Jan/2025:10:01:04] "DELETE /api/cache HTTP/1.1" 204 0
172.16.0.1 - - [01/Jan/2025:10:01:05] "GET /api/health HTTP/1.1" 200 12
192.168.1.200 - - [01/Jan/2025:10:01:06] "GET /favicon.ico HTTP/1.1" 404 0
EOF
