#!/bin/bash
age=25
score=85

if [ "$age" -gt 18 -a "$score" -gt 60 ]; then
    echo "Adult with passing score"
fi

if [ "$age" -lt 18 -o "$score" -gt 90 ]; then
    echo "Under 18 or above 90"
else
    echo "Not (under 18 or above 90)"
fi

if [ ! "$age" -eq 30 ]; then
    echo "age is not 30"
fi

if [[ "$age" -ge 18 && "$age" -le 65 ]]; then
    echo "Working age"
fi

if [[ "$score" -ge 90 || "$age" -lt 20 ]]; then
    echo "Excellent or young"
else
    echo "Not (excellent or young)"
fi
