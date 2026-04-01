#!/bin/bash
apt-get update -qq && apt-get install -y -qq mdadm > /dev/null 2>&1
rm -f /tmp/mdadm_help.txt /tmp/raid_status.txt /tmp/raid_levels.txt
