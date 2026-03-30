#!/bin/bash
mkdir -p /home/lab/mystery
echo "This is plain text" > /home/lab/mystery/image.png
echo '#!/bin/bash' > /home/lab/mystery/data.csv
echo 'echo hello' >> /home/lab/mystery/data.csv
echo '{"name": "test"}' > /home/lab/mystery/program.exe
