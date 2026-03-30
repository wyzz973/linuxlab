#!/bin/bash
sed -i 's|postgres://user:pass@localhost:5432/olddb|postgres://admin:secret123@db.production.com:5432/appdb|' challenge.txt
