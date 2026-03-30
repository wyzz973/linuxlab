#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/long_task.sh << 'EOF'
#!/bin/bash
echo "Task started at $(date)"
for i in $(seq 1 5); do
    echo "Processing step $i..."
    sleep 2
done
echo "Task completed at $(date)"
EOF
chmod +x /home/lab/long_task.sh
rm -f /home/lab/task.log
