#!/bin/bash
ssh-keygen -t ed25519 -f /tmp/test_key -N "" -q
cat /tmp/test_key.pub > /tmp/pubkey_content.txt
if [ -f /etc/ssh/sshd_config ]; then
    cp /etc/ssh/sshd_config /tmp/ssh_config.txt
    grep -q "PermitRootLogin" /tmp/ssh_config.txt || echo "PermitRootLogin prohibit-password" >> /tmp/ssh_config.txt
else
    cat > /tmp/ssh_config.txt << SSHEOF
PermitRootLogin prohibit-password
PasswordAuthentication no
PubkeyAuthentication yes
ChallengeResponseAuthentication no
UsePAM yes
X11Forwarding no
AllowTcpForwarding no
SSHEOF
fi
