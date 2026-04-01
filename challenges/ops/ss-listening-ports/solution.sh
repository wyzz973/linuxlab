#!/bin/bash
ss -tlnp > /tmp/listening.txt 2>&1
ss -tlnp | grep LISTEN | awk "{print \$4}" | rev | cut -d: -f1 | rev | sort -n | uniq > /tmp/ports.txt
