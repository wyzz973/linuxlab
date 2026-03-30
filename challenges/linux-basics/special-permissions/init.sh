#!/bin/bash
mkdir -p /home/lab
echo '#!/bin/bash' > /home/lab/special_program
echo 'echo "Running as $(whoami)"' >> /home/lab/special_program
chmod 755 /home/lab/special_program
mkdir -p /home/lab/shared_dir
chmod 777 /home/lab/shared_dir
