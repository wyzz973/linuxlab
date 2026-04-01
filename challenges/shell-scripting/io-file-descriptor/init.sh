#!/bin/bash
mkdir -p /home/learner
echo "Hello from file descriptor" > /home/learner/fd_input.txt
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
# 练习文件描述符和 exec
EOF
chmod +x /home/learner/solution.sh
