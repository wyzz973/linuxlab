#!/bin/bash
rm -f /tmp/top_ips.txt /tmp/status_codes.txt /tmp/not_found.txt
cat > /var/log/nginx_access_sim.log << LOGEOF
192.168.1.1 - - [15/Mar/2024:10:00:01 +0000] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0"
192.168.1.1 - - [15/Mar/2024:10:00:02 +0000] "GET /about HTTP/1.1" 200 1024 "-" "Mozilla/5.0"
10.0.0.5 - - [15/Mar/2024:10:00:03 +0000] "GET /api/users HTTP/1.1" 200 2048 "-" "curl/7.68"
192.168.1.1 - - [15/Mar/2024:10:00:04 +0000] "GET /missing-page HTTP/1.1" 404 162 "-" "Mozilla/5.0"
10.0.0.5 - - [15/Mar/2024:10:00:05 +0000] "POST /api/login HTTP/1.1" 200 128 "-" "curl/7.68"
172.16.0.10 - - [15/Mar/2024:10:00:06 +0000] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0"
192.168.1.1 - - [15/Mar/2024:10:00:07 +0000] "GET /admin HTTP/1.1" 403 162 "-" "Mozilla/5.0"
10.0.0.5 - - [15/Mar/2024:10:00:08 +0000] "GET /api/users HTTP/1.1" 200 2048 "-" "curl/7.68"
192.168.1.100 - - [15/Mar/2024:10:00:09 +0000] "GET /old-page HTTP/1.1" 404 162 "-" "Googlebot"
192.168.1.100 - - [15/Mar/2024:10:00:10 +0000] "GET / HTTP/1.1" 200 612 "-" "Googlebot"
10.0.0.5 - - [15/Mar/2024:10:00:11 +0000] "GET /api/data HTTP/1.1" 500 64 "-" "curl/7.68"
192.168.1.1 - - [15/Mar/2024:10:00:12 +0000] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0"
172.16.0.10 - - [15/Mar/2024:10:00:13 +0000] "GET /deleted HTTP/1.1" 404 162 "-" "Mozilla/5.0"
10.0.0.20 - - [15/Mar/2024:10:00:14 +0000] "GET / HTTP/1.1" 301 162 "-" "Mozilla/5.0"
192.168.1.1 - - [15/Mar/2024:10:00:15 +0000] "GET /contact HTTP/1.1" 200 512 "-" "Mozilla/5.0"
LOGEOF
