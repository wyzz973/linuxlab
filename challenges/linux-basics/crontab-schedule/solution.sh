#!/bin/bash
(crontab -l 2>/dev/null; echo '30 2 * * * /home/lab/backup.sh') | crontab -
