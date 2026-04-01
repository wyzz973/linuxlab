#!/bin/bash
if [ ! -f /etc/systemd/system/mydb.service ]; then
    echo "FAIL: mydb.service not found"
    exit 1
fi
if [ ! -f /etc/systemd/system/webapp.service ]; then
    echo "FAIL: webapp.service not found"
    exit 1
fi
WEBAPP=$(cat /etc/systemd/system/webapp.service)
if ! echo "$WEBAPP" | grep -q "Requires.*mydb"; then
    echo "FAIL: webapp missing Requires=mydb.service"
    exit 1
fi
if ! echo "$WEBAPP" | grep -q "After.*mydb"; then
    echo "FAIL: webapp missing After=mydb.service"
    exit 1
fi
if [ ! -f /tmp/dep_tree.txt ]; then
    echo "FAIL: /tmp/dep_tree.txt not found"
    exit 1
fi
echo "PASS"
exit 0
