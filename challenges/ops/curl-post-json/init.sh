#!/bin/bash
apt-get update -qq && apt-get install -y -qq curl python3 > /dev/null 2>&1
rm -f /tmp/get_response.txt /tmp/post_response.txt /tmp/response_headers.txt
cat > /tmp/server.py << PYEOF
from http.server import HTTPServer, BaseHTTPRequestHandler
import json

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Type", "text/plain")
        self.end_headers()
        self.wfile.write(b"Hello from GET")
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length)
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps({"received": body.decode()}).encode())
    def log_message(self, format, *args):
        pass

HTTPServer(("", 8080), Handler).serve_forever()
PYEOF
python3 /tmp/server.py &
sleep 1
