#!/bin/bash
groupadd docker 2>/dev/null || true
groupadd developers 2>/dev/null || true
userdel -r devuser 2>/dev/null || true
useradd -m -s /bin/bash -G developers devuser
