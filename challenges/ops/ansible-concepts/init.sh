#!/bin/bash
apt-get update -qq && apt-get install -y -qq python3 python3-pip > /dev/null 2>&1
pip3 install ansible > /dev/null 2>&1 || true
rm -f /tmp/ansible_version.txt /tmp/inventory.ini /tmp/setup.yml
