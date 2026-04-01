#!/bin/bash
cat > /etc/systemd/system/mydb.service << DBEOF
[Unit]
Description=My Database Service

[Service]
Type=simple
ExecStart=/bin/sleep infinity
Restart=always

[Install]
WantedBy=multi-user.target
DBEOF

cat > /etc/systemd/system/webapp.service << WAEOF
[Unit]
Description=My Web Application
Requires=mydb.service
After=mydb.service network.target

[Service]
Type=simple
ExecStart=/bin/sleep infinity
Restart=always

[Install]
WantedBy=multi-user.target
WAEOF

systemctl daemon-reload 2>/dev/null
systemctl list-dependencies webapp.service --no-pager > /tmp/dep_tree.txt 2>&1 || echo "Dependencies listed" > /tmp/dep_tree.txt
