#!/bin/bash
apt-get update -qq && apt-get install -y -qq openssh-client > /dev/null 2>&1
rm -f /tmp/test_key /tmp/test_key.pub /tmp/pubkey_content.txt /tmp/ssh_config.txt
