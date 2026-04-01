#!/bin/bash
export MY_APP="LinuxLab"
echo "$HOME"
echo "$MY_APP"
if [[ "$PATH" == *"/usr/bin"* ]]; then
    echo "PATH contains /usr/bin"
else
    echo "PATH does not contain /usr/bin"
fi
