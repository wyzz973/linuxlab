#!/bin/bash
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
text="I love apple, apple is sweet"
path="/home/user/documents/file.txt"
# 1. 替换第一个 apple 为 banana
# 2. 替换所有 apple 为 orange
# 3. 替换 path 中所有 / 为 -
EOF
chmod +x /home/learner/solution.sh
