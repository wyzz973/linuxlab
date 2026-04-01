#!/bin/bash
mkdir -p /home/learner
cat > /home/learner/users.txt << 'EOF'
alice 25 admin
bob 17 user
charlie 30 user
dave 16 admin
EOF
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
# 读取用户信息并判断权限
EOF
chmod +x /home/learner/solution.sh
