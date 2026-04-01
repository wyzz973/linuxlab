#!/bin/bash
global_var="I am global"

test_scope() {
    local local_var="I am local"
    global_var="Modified in function"
    echo "$local_var"
}

test_scope
echo "$global_var"
echo "$local_var"
