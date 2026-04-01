#!/bin/bash
text="I love apple, apple is sweet"
path="/home/user/documents/file.txt"
echo "${text/apple/banana}"
echo "${text//apple/orange}"
echo "${path//\//-}"
