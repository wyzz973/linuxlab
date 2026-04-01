#!/bin/bash
dd if=/dev/urandom of=/tmp/random.dat bs=1M count=10 2>/dev/null
dd if=/tmp/random.dat of=/tmp/random_copy.dat bs=1M 2>/dev/null
md5sum /tmp/random.dat /tmp/random_copy.dat > /tmp/dd_verify.txt
