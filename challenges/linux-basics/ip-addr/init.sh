#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq >/dev/null 2>&1 && apt-get install -y -qq --no-install-recommends iproute2 >/dev/null 2>&1
rm -f /tmp/result.txt
