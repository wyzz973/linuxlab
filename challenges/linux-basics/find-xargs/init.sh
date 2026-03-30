#!/bin/bash
mkdir -p /home/lab/project/{src,config,docs}
echo "good code" > /home/lab/project/src/main.py
echo "old code" > /home/lab/project/src/main.py.bak
echo "config" > /home/lab/project/config/app.conf
echo "old config" > /home/lab/project/config/app.conf.bak
echo "docs" > /home/lab/project/docs/readme.md
echo "old docs" > /home/lab/project/docs/readme.md.bak
echo "another backup" > /home/lab/project/old_data.bak
