#!/bin/bash
mkdir -p /home/lab/project/{src,tests,docs,config}
echo 'print("hello")' > /home/lab/project/src/main.py
echo 'import unittest' > /home/lab/project/tests/test_main.py
echo 'DB_HOST=localhost' > /home/lab/project/config/settings.conf
echo '# README' > /home/lab/project/docs/readme.md
echo 'def helper(): pass' > /home/lab/project/src/utils.py
echo '{}' > /home/lab/project/config/data.json
