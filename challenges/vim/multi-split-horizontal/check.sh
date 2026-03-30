#!/bin/bash
expected='<!DOCTYPE html>
<html>
<head>
    <title>My Page</title>
</head>
<body>
    <p>Hello World</p>
</body>
</html>'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
