#!/bin/bash
grep -q "postgres://admin:secret123@db.production.com:5432/appdb" challenge.txt
