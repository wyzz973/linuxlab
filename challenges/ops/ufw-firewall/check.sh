#!/bin/bash
if [ ! -f /tmp/ufw_status.txt ]; then
    echo "FAIL: /tmp/ufw_status.txt not found"
    exit 1
fi
if [ ! -f /tmp/ufw_setup.sh ]; then
    echo "FAIL: /tmp/ufw_setup.sh not found"
    exit 1
fi
SETUP=$(cat /tmp/ufw_setup.sh)
if ! echo "$SETUP" | grep -q "22\|ssh\|SSH"; then
    echo "FAIL: SSH rule not found"
    exit 1
fi
if ! echo "$SETUP" | grep -q "80\|http\|HTTP"; then
    echo "FAIL: HTTP rule not found"
    exit 1
fi
if ! echo "$SETUP" | grep -q "443\|https\|HTTPS"; then
    echo "FAIL: HTTPS rule not found"
    exit 1
fi
if [ ! -f /tmp/ufw_apps.txt ]; then
    echo "FAIL: /tmp/ufw_apps.txt not found"
    exit 1
fi
echo "PASS"
exit 0
