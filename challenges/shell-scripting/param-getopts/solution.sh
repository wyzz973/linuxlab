#!/bin/bash
verbose=false
name=""
age=""

while getopts "n:a:vh" opt; do
    case $opt in
        n) name="$OPTARG" ;;
        a) age="$OPTARG" ;;
        v) verbose=true ;;
        h) echo "Usage: $0 [-n name] [-a age] [-v] [-h]"; exit 0 ;;
        ?) echo "Invalid option"; exit 1 ;;
    esac
done

if [ "$verbose" = true ]; then
    echo "Verbose mode enabled"
fi
[ -n "$name" ] && echo "Name: $name"
[ -n "$age" ] && echo "Age: $age"
