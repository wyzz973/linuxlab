#!/bin/bash
rm -f /tmp/failed_count.txt /tmp/top_attacker.txt /tmp/attempted_users.txt
cat > /var/log/auth_sim.log << LOGEOF
Mar 15 08:01:23 server sshd[1234]: Failed password for root from 192.168.1.100 port 22 ssh2
Mar 15 08:01:25 server sshd[1235]: Failed password for root from 192.168.1.100 port 22 ssh2
Mar 15 08:01:27 server sshd[1236]: Failed password for admin from 192.168.1.100 port 22 ssh2
Mar 15 08:02:10 server sshd[1237]: Accepted password for ubuntu from 10.0.0.5 port 22 ssh2
Mar 15 08:03:15 server sshd[1238]: Failed password for root from 10.0.0.50 port 22 ssh2
Mar 15 08:03:18 server sshd[1239]: Failed password for test from 10.0.0.50 port 22 ssh2
Mar 15 08:04:20 server sshd[1240]: Failed password for root from 192.168.1.100 port 22 ssh2
Mar 15 08:04:22 server sshd[1241]: Failed password for admin from 192.168.1.100 port 22 ssh2
Mar 15 08:04:25 server sshd[1242]: Failed password for ubuntu from 192.168.1.100 port 22 ssh2
Mar 15 08:05:30 server sshd[1243]: Failed password for nobody from 172.16.0.10 port 22 ssh2
Mar 15 08:05:33 server sshd[1244]: Accepted password for root from 10.0.0.1 port 22 ssh2
Mar 15 08:06:40 server sshd[1245]: Failed password for root from 192.168.1.100 port 22 ssh2
Mar 15 08:06:42 server sshd[1246]: Failed password for deploy from 10.0.0.50 port 22 ssh2
LOGEOF
