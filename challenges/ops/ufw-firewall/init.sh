#!/bin/bash
# Ensure ufw is available (skip slow apt-get if already installed)
if ! command -v ufw > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq ufw > /dev/null 2>&1 || true
fi
rm -f /tmp/ufw_status.txt /tmp/ufw_setup.sh /tmp/ufw_apps.txt
