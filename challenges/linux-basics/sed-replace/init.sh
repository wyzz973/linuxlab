#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/config.txt << 'EOF'
# Application Configuration
db_host=localhost
db_port=3306
cache_host=localhost
cache_port=6379
api_url=http://localhost:8080/api
log_level=info
EOF
