#!/bin/bash
# Create the rogue process script
cat > /tmp/rogue_process.sh << 'EOF'
#!/bin/bash
while true; do sleep 10; done
EOF
chmod +x /tmp/rogue_process.sh
# Start the rogue process in background (nohup so it survives shell exit)
nohup /tmp/rogue_process.sh > /dev/null 2>&1 &
