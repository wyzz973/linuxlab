#!/bin/bash
if [ ! -f /tmp/random.dat ]; then
    echo "FAIL: /tmp/random.dat not found"
    exit 1
fi
if [ ! -f /tmp/random_copy.dat ]; then
    echo "FAIL: /tmp/random_copy.dat not found"
    exit 1
fi
if [ ! -f /tmp/dd_verify.txt ]; then
    echo "FAIL: /tmp/dd_verify.txt not found"
    exit 1
fi
MD5_ORIG=$(md5sum /tmp/random.dat | awk "{print \$1}")
MD5_COPY=$(md5sum /tmp/random_copy.dat | awk "{print \$1}")
if [ "$MD5_ORIG" = "$MD5_COPY" ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: files do not match"
exit 1
