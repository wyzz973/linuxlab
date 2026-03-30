#!/bin/bash
cd /home/lab/images && for f in *.JPG; do mv "$f" "${f%.JPG}.jpg"; done
