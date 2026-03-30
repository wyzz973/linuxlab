#!/bin/bash
sed -i '/<!-- HEADER -->/r header.html' challenge.txt
sed -i '/<!-- HEADER -->/d' challenge.txt
