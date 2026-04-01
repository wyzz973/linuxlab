#!/bin/bash
apt-get update -qq && apt-get install -y -qq coreutils > /dev/null 2>&1
rm -f /tmp/random.dat /tmp/random_copy.dat /tmp/dd_verify.txt
