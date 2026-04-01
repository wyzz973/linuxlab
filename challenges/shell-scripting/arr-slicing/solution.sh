#!/bin/bash
nums=(10 20 30 40 50 60 70 80 90 100)
echo "${nums[@]:2:3}"
nums+=(110)
echo "${#nums[@]}"
unset nums[0]
echo "${nums[@]}"
