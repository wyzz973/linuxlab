#!/bin/bash
# Create source volume with test data
docker volume create backup-source
docker run --rm -v backup-source:/data alpine sh -c '
  echo "important data" > /data/file1.txt
  echo "more data" > /data/file2.txt
  mkdir -p /data/subdir
  echo "nested" > /data/subdir/file3.txt
'
# Backup
docker run --rm -v backup-source:/source -v /tmp:/backup alpine tar czf /backup/volume-backup.tar -C /source .
# Restore to new volume
docker volume create backup-restore
docker run --rm -v backup-restore:/target -v /tmp:/backup alpine tar xzf /backup/volume-backup.tar -C /target
