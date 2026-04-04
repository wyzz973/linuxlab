#!/bin/bash
LOG=/var/log/nginx_access_sim.log

# Top 5 IPs by request count
awk '{print $1}' "$LOG" | sort | uniq -c | sort -rn | head -5 > /tmp/top_ips.txt

# Status code statistics
awk '{print $9}' "$LOG" | sort | uniq -c | sort -rn > /tmp/status_codes.txt

# All 404 request URLs
awk '$9 == "404" {print $7}' "$LOG" > /tmp/not_found.txt
