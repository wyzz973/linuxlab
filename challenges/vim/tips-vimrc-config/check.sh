#!/bin/bash
expected='" Show line numbers
set number
" Use 4 spaces for tabs
set tabstop=4
set expandtab
" Enable syntax highlighting
syntax on
" Enable auto indentation
set autoindent
" Show matching brackets
set showmatch'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
