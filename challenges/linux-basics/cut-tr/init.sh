#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/passwd_sample.txt << 'EOF'
Root:x:0:0:Root User:/root:/bin/bash
Daemon:x:1:1:Daemon:/usr/sbin:/usr/sbin/nologin
Admin:x:1000:1000:Admin User:/home/admin:/bin/bash
Guest:x:1001:1001:Guest User:/home/guest:/bin/sh
MySQL:x:27:27:MySQL Server:/var/lib/mysql:/bin/false
EOF
