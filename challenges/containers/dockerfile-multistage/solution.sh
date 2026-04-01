#!/bin/bash
mkdir -p /tmp/multistage
cat > /tmp/multistage/hello.c << 'EOF'
#include <stdio.h>
int main() { printf("Hello from multi-stage!\n"); return 0; }
EOF
cat > /tmp/multistage/Dockerfile << 'EOF'
FROM gcc:12 AS builder
WORKDIR /build
COPY hello.c .
RUN gcc -static -o hello hello.c

FROM alpine
WORKDIR /app
COPY --from=builder /build/hello .
CMD ["./hello"]
EOF
docker build -t multistage:v1 /tmp/multistage/
