#!/bin/bash
mkdir -p /home/lab/config /home/lab/scripts
echo "server_port=8080" > /home/lab/config/app.conf
echo "db_host=localhost" >> /home/lab/config/app.conf
echo '#!/bin/bash' > /home/lab/scripts/start.sh
echo 'echo "Starting..."' >> /home/lab/scripts/start.sh
rm -f /home/lab/app.conf /home/lab/bin
