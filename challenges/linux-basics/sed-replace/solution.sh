#!/bin/bash
sed 's/localhost/db.production.com/g' /home/lab/config.txt > /tmp/result.txt
