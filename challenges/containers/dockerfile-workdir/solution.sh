#!/bin/bash
mkdir -p /tmp/workdir-app
echo 'echo "app is running"' > /tmp/workdir-app/app.sh
cat > /tmp/workdir-app/Dockerfile << 'EOF'
FROM alpine
WORKDIR /app
COPY app.sh .
RUN chmod +x app.sh
CMD ["./app.sh"]
EOF
docker build -t workdir-app:v1 /tmp/workdir-app/
