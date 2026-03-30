#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/server.log << 'EOF'
2025-01-01 00:00:01 Server started
2025-01-01 00:00:02 Loading configuration
2025-01-01 00:00:03 Database connected
2025-01-01 00:00:04 Cache initialized
2025-01-01 00:00:05 Listening on port 8080
2025-01-01 00:01:00 GET /api/users 200
2025-01-01 00:01:01 GET /api/products 200
2025-01-01 00:02:00 POST /api/orders 201
2025-01-01 00:03:00 GET /api/users/1 200
2025-01-01 00:04:00 ERROR: Connection timeout
2025-01-01 00:04:01 Retrying connection...
2025-01-01 00:04:02 Connection restored
EOF
