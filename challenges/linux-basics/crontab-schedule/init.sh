#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq >/dev/null 2>&1 && apt-get install -y -qq --no-install-recommends cron >/dev/null 2>&1
mkdir -p /home/lab
echo '#!/bin/bash' > /home/lab/backup.sh
echo 'tar czf /backup/data-$(date +%Y%m%d).tar.gz /home/lab/data' >> /home/lab/backup.sh
chmod +x /home/lab/backup.sh
crontab -r 2>/dev/null || true
