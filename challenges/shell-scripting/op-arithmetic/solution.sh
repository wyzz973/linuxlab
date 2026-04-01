#!/bin/bash
a=10
b=3
echo $(expr $a + $b)
echo $(expr $a - $b)
echo $(expr $a \* $b)
echo $(expr $a / $b)
echo $(expr $a % $b)
echo $(( (a + b) * 2 ))
