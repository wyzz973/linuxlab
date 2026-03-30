#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/config.txt << 'EOF'
[server]
host=localhost
port=8080

[database]
host=db.local
port=3306

[cache]
enabled=true
EOF
