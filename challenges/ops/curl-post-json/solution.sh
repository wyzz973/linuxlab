#!/bin/bash
curl -s http://localhost:8080/ > /tmp/get_response.txt
curl -s -X POST -H "Content-Type: application/json" -d "{\"name\":\"test\"}" http://localhost:8080/ > /tmp/post_response.txt
curl -sI http://localhost:8080/ > /tmp/response_headers.txt
