#!/bin/bash
awk '{gsub(/^- /, NR". "); print}' challenge.txt > /tmp/vim_macro.txt && mv /tmp/vim_macro.txt challenge.txt
