#!/bin/bash
mkdir -p /home/lab/scripts
echo '#!/bin/bash' > /home/lab/scripts/deploy.sh
echo 'echo deploy' >> /home/lab/scripts/deploy.sh
echo '#!/bin/bash' > /home/lab/scripts/backup.sh
echo 'echo backup' >> /home/lab/scripts/backup.sh
echo '#!/bin/bash' > /home/lab/scripts/cleanup.sh
echo 'echo cleanup' >> /home/lab/scripts/cleanup.sh
echo 'config=true' > /home/lab/scripts/settings.conf
chmod 644 /home/lab/scripts/*.sh
chmod 644 /home/lab/scripts/settings.conf
