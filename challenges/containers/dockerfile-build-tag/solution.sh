#!/bin/bash
mkdir -p /tmp/tagdemo
cat > /tmp/tagdemo/Dockerfile << 'EOF'
FROM alpine
CMD ["echo", "tag demo"]
EOF
docker build -t tagdemo:v1 /tmp/tagdemo/
docker tag tagdemo:v1 tagdemo:latest
docker tag tagdemo:v1 myregistry/tagdemo:v1
docker images tagdemo > /tmp/tag-list.txt
docker images myregistry/tagdemo >> /tmp/tag-list.txt
