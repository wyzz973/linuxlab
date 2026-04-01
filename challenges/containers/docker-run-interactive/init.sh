#!/bin/bash
docker pull ubuntu:latest 2>/dev/null
docker rm -f interactive-ubuntu 2>/dev/null
