#!/bin/bash
expected='# Task List
- [x] Setup project structure
- [ ] DONE: Implement login feature
- [x] Write unit tests
- [x] Deploy to staging'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
