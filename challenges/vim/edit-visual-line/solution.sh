#!/bin/bash
sed -i 's/^\([a-z_]*\)=/\U\1=/' challenge.txt
