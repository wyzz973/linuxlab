#!/bin/bash
rm -f /tmp/result.txt
# Make sure there's at least one listening port
python3 -c "
import socket, time, threading
s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('0.0.0.0', 9999))
s.listen(1)
time.sleep(300)
" &>/dev/null &
sleep 1
