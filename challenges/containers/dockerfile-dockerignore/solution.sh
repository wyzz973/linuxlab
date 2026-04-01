#!/bin/bash
mkdir -p /tmp/ignore-demo/logs
echo 'echo "production app"' > /tmp/ignore-demo/app.sh
echo 'echo "test"' > /tmp/ignore-demo/test.sh
echo '# README' > /tmp/ignore-demo/README.md
echo 'debug info' > /tmp/ignore-demo/logs/debug.log
cat > /tmp/ignore-demo/.dockerignore << 'EOF'
*.md
test.sh
logs/
EOF
cat > /tmp/ignore-demo/Dockerfile << 'EOF'
FROM alpine
WORKDIR /app
COPY . .
CMD ["sh", "app.sh"]
EOF
docker build -t ignore-demo:v1 /tmp/ignore-demo/
