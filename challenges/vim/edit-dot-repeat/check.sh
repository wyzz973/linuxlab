#!/bin/bash
expected='let name = "Alice";
let age = 30;
let city = "Shanghai";
let role = "engineer";
let level = "senior";'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
