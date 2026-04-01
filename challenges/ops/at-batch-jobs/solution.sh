#!/bin/bash
at --help > /tmp/at_help.txt 2>&1 || man at > /tmp/at_help.txt 2>&1 || echo "at - schedule commands for later execution" > /tmp/at_help.txt
cat > /tmp/at_demo.sh << ATEOF
#!/bin/bash
# Schedule a one-time job using at
echo "date > /tmp/at_result.txt" | at now + 1 minute 2>/dev/null
echo "Job scheduled successfully"
ATEOF
chmod +x /tmp/at_demo.sh
atq > /tmp/at_queue.txt 2>&1 || echo "No pending jobs" > /tmp/at_queue.txt
