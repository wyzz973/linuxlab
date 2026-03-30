#!/bin/bash
! grep -q "TODO" challenge.txt && grep -c "DONE" challenge.txt | grep -q "5"
