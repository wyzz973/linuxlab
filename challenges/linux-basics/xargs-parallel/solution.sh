#!/bin/bash
cat /home/lab/urls.txt | xargs -I{} bash -c 'touch /home/lab/downloads/$(basename "{}")'
