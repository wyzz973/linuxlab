#!/bin/bash
url="https://www.runoob.com/linux/linux-shell.html"
echo ${#url}
echo ${url:12:6}
expr index "$url" l
