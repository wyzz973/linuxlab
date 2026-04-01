#!/bin/bash
cat > /usr/local/bin/myapp.sh << APPEOF
#!/bin/bash
while true; do
    echo "\$(date)" >> /var/log/myapp.log
    sleep 1
done
APPEOF
chmod +x /usr/local/bin/myapp.sh
cat > /etc/systemd/system/myapp.service << SVCEOF
[Unit]
Description=My Custom App
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/myapp.sh
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SVCEOF
systemctl daemon-reload 2>/dev/null
