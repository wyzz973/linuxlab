#!/bin/bash
mkdir -p /home/learner
cat > /home/learner/calculations.txt << 'EOF'
10 + 5
20 - 8
6 * 7
100 / 4
17 % 3
EOF
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
# 读取 calculations.txt 并计算每行结果
EOF
chmod +x /home/learner/solution.sh
