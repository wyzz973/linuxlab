#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq >/dev/null 2>&1 && apt-get install -y -qq --no-install-recommends file >/dev/null 2>&1
mkdir -p /home/lab/mystery
echo "This is plain text" > /home/lab/mystery/image.png
echo '#!/bin/bash' > /home/lab/mystery/data.csv
echo 'echo hello' >> /home/lab/mystery/data.csv
echo '{"name": "test"}' > /home/lab/mystery/program.exe
