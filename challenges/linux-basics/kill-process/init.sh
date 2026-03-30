#!/bin/bash
pkill -f rogue_process 2>/dev/null || true
sleep 0.5
# Start a background process that simulates a rogue process
nohup bash -c 'while true; do sleep 10; done' > /dev/null 2>&1 &
# Rename it by creating a wrapper
cat > /tmp/rogue_process.sh << 'EOF'
#!/bin/bash
while true; do sleep 10; done
EOF
chmod +x /tmp/rogue_process.sh
nohup /tmp/rogue_process.sh > /dev/null 2>&1 &
