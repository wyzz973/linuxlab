#!/bin/bash
docker network rm my-network isolated-net 2>/dev/null
rm -f /tmp/network-list.txt /tmp/network-list-after.txt
