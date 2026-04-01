#!/bin/bash
cat > /tmp/my_crontab << CRONEOF
# 每天凌晨 2 点备份数据库
0 2 * * * /usr/local/bin/backup_db.sh >> /var/log/backup_db.sh.log 2>&1
# 每 5 分钟检查磁盘空间
*/5 * * * * /usr/local/bin/check_disk.sh >> /var/log/check_disk.sh.log 2>&1
# 每周一上午 9 点发送周报
0 9 * * 1 /usr/local/bin/weekly_report.sh >> /var/log/weekly_report.sh.log 2>&1
# 每月 1 号清理日志
0 0 1 * * /usr/local/bin/cleanup_logs.sh >> /var/log/cleanup_logs.sh.log 2>&1
CRONEOF
crontab -l > /tmp/current_cron.txt 2>&1 || echo "No crontab for current user" > /tmp/current_cron.txt
