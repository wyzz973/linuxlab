#!/bin/bash
mkdir -p /home/learner
cat > /home/learner/words.txt << 'EOF'
hello
world
shell
EOF
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
# 读取 /home/learner/words.txt 的每一行
# 输出 "字符串: 长度"
EOF
chmod +x /home/learner/solution.sh
