#!/bin/bash
expected='event1: 15-01-2024, Conference
event2: 22-03-2024, Workshop
event3: 10-06-2024, Meetup
event4: 05-09-2024, Hackathon
event5: 31-12-2024, New Year'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
