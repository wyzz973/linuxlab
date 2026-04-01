#!/bin/bash
mkdir -p /tmp/greeter
cat > /tmp/greeter/Dockerfile << 'EOF'
FROM alpine
ENTRYPOINT ["echo", "Hello"]
CMD ["World"]
EOF
docker build -t greeter:v1 /tmp/greeter/
