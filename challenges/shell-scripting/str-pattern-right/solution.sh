#!/bin/bash
filepath="/home/user/documents/report.tar.gz"
echo "${filepath%.*}"
echo "${filepath%%.*}"
echo "${filepath%/*}"
