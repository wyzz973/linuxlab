#!/bin/bash
userdel -r operator 2>/dev/null || true
useradd -m operator
# Ensure sudo group exists
groupadd sudo 2>/dev/null || true
