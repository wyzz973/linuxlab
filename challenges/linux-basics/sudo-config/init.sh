#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq >/dev/null 2>&1 && apt-get install -y -qq --no-install-recommends sudo >/dev/null 2>&1
userdel -r operator 2>/dev/null || true
groupdel operator 2>/dev/null || true
useradd -m operator
# Ensure sudo group exists
groupadd sudo 2>/dev/null || true
