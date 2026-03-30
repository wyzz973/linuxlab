#!/bin/bash
apt update -y 2>/dev/null
apt install -y tree 2>/dev/null
dpkg -l tree > /tmp/result.txt 2>/dev/null || apt list --installed 2>/dev/null | grep tree > /tmp/result.txt
