#!/bin/bash
# Ensure ssh-keygen is available (skip slow apt-get if already installed)
if ! command -v ssh-keygen > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq openssh-client > /dev/null 2>&1 || true
fi
rm -f /tmp/test_key /tmp/test_key.pub /tmp/pubkey_content.txt /tmp/ssh_config.txt
