#!/bin/bash
mkdir -p /home/learner
cat > /home/learner/items.txt << 'EOF'
apple
banana
cherry
EOF
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
# 练习 for 循环遍历列表
EOF
chmod +x /home/learner/solution.sh
