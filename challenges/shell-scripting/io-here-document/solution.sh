#!/bin/bash
cat > /home/learner/config.txt <<EOF
[server]
host=localhost
port=8080
debug=true
EOF

name="World"
cat <<EOF
Hello, $name! Today is $(date +%Y-%m-%d)
EOF

cat <<'EOF'
Price: $100
EOF

cat /home/learner/config.txt
