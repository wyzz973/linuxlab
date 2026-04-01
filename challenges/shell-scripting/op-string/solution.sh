#!/bin/bash
str1="hello"
str2="world"
str3=""
[ "$str1" = "$str2" ] && echo "equal" || echo "not equal"
[ "$str1" != "$str2" ] && echo "different" || echo "same"
[ -z "$str3" ] && echo "str3 is empty" || echo "str3 is not empty"
[ -n "$str1" ] && echo "str1 is not empty" || echo "str1 is empty"
[ "$str1" ] && echo "str1 has value" || echo "str1 has no value"
