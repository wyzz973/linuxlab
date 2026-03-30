#!/bin/bash
if ! id devuser &>/dev/null; then
    echo "FAIL: User devuser not found"
    exit 1
fi
groups=$(groups devuser 2>/dev/null)
shell=$(getent passwd devuser | cut -d: -f7)
if echo "$groups" | grep -q "docker" && echo "$groups" | grep -q "developers" && [ "$shell" = "/bin/zsh" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Groups=$groups, Shell=$shell"
    exit 1
fi
