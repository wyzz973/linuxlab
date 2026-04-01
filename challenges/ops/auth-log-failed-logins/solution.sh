#!/bin/bash
grep -c "Failed password" /var/log/auth_sim.log > /tmp/failed_count.txt
grep "Failed password" /var/log/auth_sim.log | awk "{print \$(NF-3)}" | sort | uniq -c | sort -rn | head -1 | awk "{print \$2}" > /tmp/top_attacker.txt
grep "Failed password" /var/log/auth_sim.log | awk "{for(i=1;i<=NF;i++) if(\$i==\"for\") print \$(i+1)}" | sort -u > /tmp/attempted_users.txt
