#!/bin/bash
groupadd webteam
groupadd -g 2000 backend
groupmod -n frontend webteam
