#!/bin/bash
rm -f /tmp/broken_status.txt /tmp/broken_logs.txt /tmp/fixed_service.txt
# Create a broken service (wrong path)
cat > /etc/systemd/system/broken.service << BRKEOF
[Unit]
Description=Broken Service for Troubleshooting

[Service]
Type=simple
ExecStart=/usr/local/bin/nonexistent_app.sh
Restart=always

[Install]
WantedBy=multi-user.target
BRKEOF
# Create the actual script at the correct path
cat > /usr/local/bin/fixed_app.sh << FIXEOF
#!/bin/bash
while true; do
    echo "Running..." >> /var/log/fixed_app.log
    sleep 5
done
FIXEOF
chmod +x /usr/local/bin/fixed_app.sh
systemctl daemon-reload 2>/dev/null
