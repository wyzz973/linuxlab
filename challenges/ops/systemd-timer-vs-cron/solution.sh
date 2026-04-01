#!/bin/bash
cat > /tmp/timer_vs_cron.txt << CMPEOF
systemd timer vs cron 对比：

systemd timer 优势：
- 集成 journalctl 日志系统，方便查看执行历史
- 支持服务依赖关系（After, Requires）
- Persistent=true 可以补偿错过的执行
- 精确到微秒级的调度精度
- 可以用 systemctl 统一管理

cron 优势：
- 语法简单直观，学习成本低
- 历史悠久，文档丰富
- 配置文件集中管理
- 适合简单的定时任务
- 不需要创建多个文件

建议：简单任务用 cron，需要依赖管理和高级特性用 systemd timer
CMPEOF

cat > /etc/systemd/system/backup.service << BSEOF
[Unit]
Description=Daily Backup

[Service]
Type=oneshot
ExecStart=/usr/local/bin/backup.sh
BSEOF

cat > /etc/systemd/system/backup.timer << BTEOF
[Unit]
Description=Run backup daily at 3am

[Timer]
OnCalendar=*-*-* 03:00:00
Persistent=true

[Install]
WantedBy=timers.target
BTEOF

systemctl daemon-reload 2>/dev/null
systemctl list-timers --all --no-pager > /tmp/all_timers.txt 2>&1
