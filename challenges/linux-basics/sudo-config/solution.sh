#!/bin/bash
usermod -aG sudo operator 2>/dev/null || usermod -aG wheel operator
