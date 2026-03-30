#!/bin/bash
mkdir -p /home/lab
echo '#!/bin/bash' > /home/lab/deploy.sh
echo 'echo "Deploying application..."' >> /home/lab/deploy.sh
echo 'cp -r /app/* /var/www/' >> /home/lab/deploy.sh
chmod 644 /home/lab/deploy.sh
