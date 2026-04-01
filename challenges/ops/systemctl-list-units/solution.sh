#!/bin/bash
systemctl list-units --type=service --no-pager > /tmp/loaded_services.txt 2>&1
systemctl --failed --no-pager > /tmp/failed_units.txt 2>&1
systemctl list-unit-files --no-pager > /tmp/unit_files.txt 2>&1
