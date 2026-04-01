#!/bin/bash
rm -f /tmp/timer_vs_cron.txt /tmp/all_timers.txt
rm -f /etc/systemd/system/backup.service /etc/systemd/system/backup.timer
mkdir -p /usr/local/bin
echo "#!/bin/bash" > /usr/local/bin/backup.sh
echo "echo backup done" >> /usr/local/bin/backup.sh
chmod +x /usr/local/bin/backup.sh
