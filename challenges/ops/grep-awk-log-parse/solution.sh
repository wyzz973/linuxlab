#!/bin/bash
grep "ERROR" /var/log/app_sim.log > /tmp/errors.txt
grep "ERROR" /var/log/app_sim.log | awk "{print \$1, \$2, \$NF}" > /tmp/error_summary.txt
grep "ERROR" /var/log/app_sim.log | awk "{split(\$2,a,\":\"); print a[1]\":00\"}" | sort | uniq -c > /tmp/hourly_errors.txt
