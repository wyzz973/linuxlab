#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/nginx.conf << 'EOF'
server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend:3000;
    }
    
    location /api {
        proxy_pass http://api:5000;
    }
}
EOF
