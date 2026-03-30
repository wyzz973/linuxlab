#!/bin/bash
apt remove -y nano 2>/dev/null || true
apt autoremove -y 2>/dev/null || true
dpkg -l > /tmp/result.txt 2>/dev/null || echo "done" > /tmp/result.txt
