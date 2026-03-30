#!/bin/bash
for f in /home/lab/webroot /home/lab/webroot/index.html /home/lab/webroot/css/style.css /home/lab/webroot/js/app.js; do
    perms=$(stat -c %a "$f" 2>/dev/null || stat -f %Lp "$f" 2>/dev/null)
    if [ "$perms" != "755" ]; then
        echo "FAIL: $f has permissions $perms, expected 755"
        exit 1
    fi
done
echo "PASS"
exit 0
