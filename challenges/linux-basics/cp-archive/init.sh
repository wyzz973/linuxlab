#!/bin/bash
mkdir -p /home/lab/production/{app,config,logs}
echo "const express = require('express');" > /home/lab/production/app/server.js
echo "PORT=3000" > /home/lab/production/config/env
echo "2025-01-01 app started" > /home/lab/production/logs/app.log
ln -sf /home/lab/production/config/env /home/lab/production/app/.env
rm -rf /home/lab/backup_production
