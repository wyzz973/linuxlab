#!/bin/bash
ufw status verbose > /tmp/ufw_status.txt 2>&1
cat > /tmp/ufw_setup.sh << UFWEOF
#!/bin/bash
# UFW Firewall Setup Script
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw --force enable
UFWEOF
chmod +x /tmp/ufw_setup.sh
ufw app list > /tmp/ufw_apps.txt 2>&1 || echo "No app profiles" > /tmp/ufw_apps.txt
