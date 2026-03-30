#!/bin/bash
prog_perms=$(stat -c %a /home/lab/special_program 2>/dev/null || stat -f %Lp /home/lab/special_program 2>/dev/null)
dir_perms=$(stat -c %a /home/lab/shared_dir 2>/dev/null || stat -f %Lp /home/lab/shared_dir 2>/dev/null)
if [ "$prog_perms" = "4755" ] && [ "$dir_perms" = "1777" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Program=$prog_perms (expected 4755), Dir=$dir_perms (expected 1777)"
    exit 1
fi
