#!/bin/bash
sed -i 's/FILL_1/set number/' challenge.txt
sed -i 's/FILL_2/set tabstop=4/' challenge.txt
sed -i 's/FILL_3/set expandtab/' challenge.txt
sed -i 's/FILL_4/syntax on/' challenge.txt
sed -i 's/FILL_5/set autoindent/' challenge.txt
sed -i 's/FILL_6/set showmatch/' challenge.txt
