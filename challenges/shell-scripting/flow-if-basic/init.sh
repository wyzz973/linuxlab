#!/bin/bash
echo "42" > /home/learner/number.txt
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
num=$(cat /home/learner/number.txt)
# 判断 num 是正数、负数还是零
EOF
chmod +x /home/learner/solution.sh
