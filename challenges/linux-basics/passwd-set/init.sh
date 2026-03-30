#!/bin/bash
userdel -r testuser 2>/dev/null || true
useradd -m testuser
