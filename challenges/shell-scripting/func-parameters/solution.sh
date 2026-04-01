#!/bin/bash
user_info() {
    echo "Name: $1, Age: $2, Job: $3"
}

count_args() {
    echo "Argument count: $#"
    echo "All arguments: $@"
}

user_info "Alice" 28 "Engineer"
count_args a b c d e
