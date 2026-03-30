#!/bin/bash
expected='DATABASE_HOST=localhost
DATABASE_PORT=5432
APP_SECRET=mykey123
LOG_LEVEL=info'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
