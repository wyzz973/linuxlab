#!/bin/bash
for f in deploy.sh backup.sh cleanup.sh; do
    perms=$(stat -c %a /home/lab/scripts/$f 2>/dev/null || stat -f %Lp /home/lab/scripts/$f 2>/dev/null)
    if ! echo "$perms" | grep -q "[0-9]*[1357][0-9]*"; then
        # Check if execute bit is set
        if [ ! -x /home/lab/scripts/$f ]; then
            echo "FAIL: /home/lab/scripts/$f is not executable"
            exit 1
        fi
    fi
done
# settings.conf should NOT have execute permission
if [ -x /home/lab/scripts/settings.conf ]; then
    echo "FAIL: settings.conf should not be executable (it's not a .sh file)"
    exit 1
fi
echo "PASS"
exit 0
