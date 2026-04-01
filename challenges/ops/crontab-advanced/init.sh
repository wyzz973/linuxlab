#!/bin/bash
apt-get update -qq && apt-get install -y -qq cron > /dev/null 2>&1
rm -f /tmp/my_crontab /tmp/current_cron.txt
# Create dummy scripts
for script in backup_db.sh check_disk.sh weekly_report.sh cleanup_logs.sh; do
    echo "#!/bin/bash" > "/usr/local/bin/$script"
    echo "echo \"Running $script at \$(date)\" >> /var/log/$script.log" >> "/usr/local/bin/$script"
    chmod +x "/usr/local/bin/$script"
done
