#!/bin/bash
mkdir -p /tmp/labeled-app
cat > /tmp/labeled-app/Dockerfile << 'EOF'
FROM alpine
ARG APP_VERSION=1.0
LABEL maintainer="student@example.com"
LABEL version="${APP_VERSION}"
CMD ["echo", "labeled app"]
EOF
docker build --build-arg APP_VERSION=2.0 -t labeled-app:v2 /tmp/labeled-app/
