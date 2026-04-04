#!/bin/bash
# Clean up any previous attempt
rm -f /tmp/get_response.txt /tmp/post_response.txt /tmp/response_headers.txt

# Create a hint file so the user knows the expected scenario
cat > /tmp/README_challenge.txt << 'EOF'
Challenge: curl POST JSON

Imagine a local HTTP server is running on port 8080.
Write the correct curl commands and save example output to the required files.

1. /tmp/get_response.txt  - should contain a plausible GET response (e.g. "Hello from GET")
2. /tmp/post_response.txt - should contain a plausible JSON POST response (e.g. {"received":"{\"name\":\"test\"}"})
3. /tmp/response_headers.txt - should contain HTTP response headers (must include "HTTP/" line)
EOF
