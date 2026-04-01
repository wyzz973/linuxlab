#!/bin/bash
greet() {
    echo "Hello from function!"
}

function show_date {
    echo "$(date +%Y-%m-%d)"
}

add() {
    echo $(($1 + $2))
}

greet
show_date
add 10 20
