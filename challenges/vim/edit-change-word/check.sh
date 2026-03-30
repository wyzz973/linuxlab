#!/bin/bash
expected='def get_user_name(user):
    return user.name

def set_password(user, pwd):
    user.password = pwd

def is_logged_in(user):
    return user.active'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
