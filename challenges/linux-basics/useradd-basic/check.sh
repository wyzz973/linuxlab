#!/bin/bash
if id developer &>/dev/null; then
    shell=$(getent passwd developer | cut -d: -f7)
    home=$(getent passwd developer | cut -d: -f6)
    if [ "$shell" = "/bin/bash" ] && [ -d "$home" ]; then
        echo "PASS"
        exit 0
    else
        echo "FAIL: Shell=$shell, Home exists=$([ -d $home ] && echo yes || echo no)"
        exit 1
    fi
else
    echo "FAIL: User developer does not exist"
    exit 1
fi
