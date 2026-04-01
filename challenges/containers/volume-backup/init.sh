#!/bin/bash
docker pull alpine:latest 2>/dev/null
docker volume rm backup-source backup-restore 2>/dev/null
rm -f /tmp/volume-backup.tar
