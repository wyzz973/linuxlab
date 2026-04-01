#!/bin/bash
mkdir -p /tmp/envapp
cat > /tmp/envapp/Dockerfile << 'EOF'
FROM alpine
ENV APP_NAME=MyApp
ENV APP_VERSION=1.0
EXPOSE 8080
CMD ["sh", "-c", "echo $APP_NAME $APP_VERSION && sleep 3600"]
EOF
docker build -t envapp:v1 /tmp/envapp/
docker run -d --name env-test envapp:v1
