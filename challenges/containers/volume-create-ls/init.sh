#!/bin/bash
docker volume rm app-data db-data 2>/dev/null
rm -f /tmp/volume-info.txt /tmp/volume-list.txt
