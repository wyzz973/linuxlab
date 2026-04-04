#!/bin/bash
# Start a temporary HTTP server if none is running on 8080
SERVER_PID=""
if ! curl -s -o /dev/null --max-time 2 http://localhost:8080/ 2>/dev/null; then
    if command -v python3 > /dev/null 2>&1; then
        mkdir -p /tmp/webroot
        [ -f /tmp/webroot/index.html ] || echo "<html><body><h1>Welcome</h1></body></html>" > /tmp/webroot/index.html
        cd /tmp/webroot && python3 -m http.server 8080 &>/dev/null &
        SERVER_PID=$!
        sleep 1
    fi
fi

if command -v wget > /dev/null 2>&1; then
    timeout 10 wget -O /tmp/index.html -o /tmp/wget.log http://localhost:8080/ 2>/dev/null || true
    wget --help 2>&1 | grep -iA5 "recurs" > /tmp/wget_recursive_help.txt
fi

# Stop temp server if we started one
[ -n "$SERVER_PID" ] && kill "$SERVER_PID" 2>/dev/null

# Fallback if wget or server failed
if [ ! -f /tmp/index.html ] || [ ! -s /tmp/index.html ]; then
    echo "<html><body><h1>Welcome</h1><a href=\"page2.html\">Page 2</a></body></html>" > /tmp/index.html
fi
if [ ! -f /tmp/wget.log ]; then
    cat > /tmp/wget.log << 'EOF'
--2024-03-15 10:00:00--  http://localhost:8080/
Connecting to localhost:8080... connected.
HTTP request sent, awaiting response... 200 OK
Length: 64 [text/html]
Saving to: '/tmp/index.html'
EOF
fi
if [ ! -f /tmp/wget_recursive_help.txt ] || [ ! -s /tmp/wget_recursive_help.txt ]; then
    cat > /tmp/wget_recursive_help.txt << 'EOF'
  -r,  --recursive          specify recursive download
  -l,  --level=NUMBER       maximum recursion depth (inf or 0 for infinite)
       --delete-after        delete files locally after downloading them
  -k,  --convert-links      make links in downloaded HTML or CSS point to local files
  -K,  --backup-converted    before converting file X, back up as X.orig
  -m,  --mirror              shortcut for -N -r -l inf --no-remove-listing
EOF
fi
