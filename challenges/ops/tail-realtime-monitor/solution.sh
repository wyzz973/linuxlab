#!/bin/bash
tail -n 20 /var/log/live_sim.log > /tmp/tail_output.txt
timeout 3 tail -f /var/log/live_sim.log 2>/dev/null | grep ERROR > /tmp/live_errors.txt &
sleep 4
wc -l < /var/log/live_sim.log | tr -d " " > /tmp/line_count.txt
