#!/bin/bash
export MY_APP_ENV=production
env > /tmp/env_result.txt
echo $PATH > /tmp/path_result.txt
