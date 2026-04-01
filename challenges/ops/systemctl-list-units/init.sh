#!/bin/bash
apt-get update -qq && apt-get install -y -qq systemd > /dev/null 2>&1
rm -f /tmp/loaded_services.txt /tmp/failed_units.txt /tmp/unit_files.txt
