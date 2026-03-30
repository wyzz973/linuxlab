#!/bin/bash
expected='server {
    listen 80;
    server_name example.com;
    location / {
        proxy_pass http://backend;
    }
}'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
