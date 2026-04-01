#!/bin/bash
openssl req -x509 -newkey rsa:2048 -keyout /tmp/server.key -out /tmp/server.crt -days 365 -nodes -subj "/CN=localhost/O=LinuxLab/C=CN" 2>/dev/null
openssl x509 -in /tmp/server.crt -text -noout > /tmp/cert_info.txt
CERT_MD5=$(openssl x509 -noout -modulus -in /tmp/server.crt | openssl md5)
KEY_MD5=$(openssl rsa -noout -modulus -in /tmp/server.key | openssl md5)
if [ "$CERT_MD5" = "$KEY_MD5" ]; then
    echo "MATCH: Certificate and key match" > /tmp/cert_verify.txt
    echo "Cert modulus MD5: $CERT_MD5" >> /tmp/cert_verify.txt
    echo "Key modulus MD5: $KEY_MD5" >> /tmp/cert_verify.txt
else
    echo "MISMATCH: Certificate and key do not match" > /tmp/cert_verify.txt
fi
