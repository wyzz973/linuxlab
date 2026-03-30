#!/bin/bash
expected='api_url = "https://api.example.com/v1"
cdn_url = "https://cdn.example.com/assets"
auth_url = "https://auth.example.com/login"
docs_url = "https://docs.example.com/api"
webhook = "https://hooks.example.com/notify"'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
