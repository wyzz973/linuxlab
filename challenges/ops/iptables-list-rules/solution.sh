#!/bin/bash
# Try real iptables first
if iptables -L -n --line-numbers > /tmp/iptables_rules.txt 2>&1 && \
   grep -qi "chain\|target\|ACCEPT\|DROP" /tmp/iptables_rules.txt 2>/dev/null; then
    iptables -t nat -L -n > /tmp/iptables_nat.txt 2>&1
    iptables -L -n | grep -cE "^[0-9]+" > /tmp/rule_count.txt 2>&1 || echo "0" > /tmp/rule_count.txt
else
    # Fallback: create simulated output for environments without NET_ADMIN
    cat > /tmp/iptables_rules.txt << 'EOF'
Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination
1    ACCEPT     tcp  --  0.0.0.0/0            0.0.0.0/0            tcp dpt:22
2    ACCEPT     tcp  --  0.0.0.0/0            0.0.0.0/0            tcp dpt:80
3    ACCEPT     tcp  --  0.0.0.0/0            0.0.0.0/0            tcp dpt:443

Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
num  target     prot opt source               destination
EOF

    cat > /tmp/iptables_nat.txt << 'EOF'
Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination

Chain INPUT (policy ACCEPT)
target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination

Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination
EOF

    echo "3" > /tmp/rule_count.txt
fi
