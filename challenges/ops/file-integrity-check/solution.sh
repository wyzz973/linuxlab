#!/bin/bash
sha256sum /etc/passwd > /tmp/passwd_sha256.txt
md5sum /etc/passwd > /tmp/passwd_md5.txt
find /etc -maxdepth 2 -name "*.conf" -exec sha256sum {} \; > /tmp/conf_checksums.txt 2>/dev/null
