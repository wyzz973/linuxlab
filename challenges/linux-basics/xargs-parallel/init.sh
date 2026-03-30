#!/bin/bash
mkdir -p /home/lab/downloads
cat > /home/lab/urls.txt << 'EOF'
http://example.com/files/report.pdf
http://example.com/files/data.csv
http://example.com/images/photo.jpg
http://example.com/docs/manual.txt
EOF
rm -f /home/lab/downloads/*
