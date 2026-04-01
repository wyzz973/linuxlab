#!/bin/bash
fruits=(apple banana cherry date elderberry)
echo "${fruits[0]}"
echo "${fruits[2]}"
echo "${fruits[@]}"
fruits[1]="blueberry"
echo "${fruits[@]}"
