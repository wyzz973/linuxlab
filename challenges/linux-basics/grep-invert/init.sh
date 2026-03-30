#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/processes.txt << 'EOF'
root      1  0.0  0.1  init
root      2  0.0  0.0  [kthreadd]
root      3  0.0  0.0  [ksoftirqd/0]
www       100  0.5  1.2  nginx: master process
www       101  0.3  0.8  nginx: worker process
mysql     200  2.0  5.0  /usr/sbin/mysqld
root      300  0.0  0.1  /usr/sbin/sshd
root      400  0.0  0.0  grep nginx
root      401  0.0  0.0  grep mysql
EOF
