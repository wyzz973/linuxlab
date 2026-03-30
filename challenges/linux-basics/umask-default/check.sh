#!/bin/bash
if [ ! -f /home/lab/secure_file.txt ]; then
    echo "FAIL: /home/lab/secure_file.txt not found"
    exit 1
fi
if [ ! -d /home/lab/secure_dir ]; then
    echo "FAIL: /home/lab/secure_dir not found"
    exit 1
fi
file_perms=$(stat -c %a /home/lab/secure_file.txt 2>/dev/null || stat -f %Lp /home/lab/secure_file.txt 2>/dev/null)
dir_perms=$(stat -c %a /home/lab/secure_dir 2>/dev/null || stat -f %Lp /home/lab/secure_dir 2>/dev/null)
if [ "$file_perms" = "640" ] && [ "$dir_perms" = "750" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: File=$file_perms (expected 640), Dir=$dir_perms (expected 750)"
    exit 1
fi
