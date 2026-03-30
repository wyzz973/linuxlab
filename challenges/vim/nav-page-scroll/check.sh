#!/bin/bash
grep -q "FIXED: Disk write failed" challenge.txt && ! grep -q "ERROR" challenge.txt
