#!/bin/bash
mkdir -p /home/learner/testdir/subdir
echo "Hello World" > /home/learner/testdir/regular.txt
touch /home/learner/testdir/empty.txt
echo '#!/bin/bash' > /home/learner/testdir/script.sh
chmod +x /home/learner/testdir/script.sh
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
# 使用文件测试运算符检测文件
base="/home/learner/testdir"
EOF
chmod +x /home/learner/solution.sh
