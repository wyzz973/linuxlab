#!/bin/bash
mkdir -p /home/learner
cat > /home/learner/files.txt << 'EOF'
setup.sh
readme.md
photo.jpg
archive.tar.gz
data.csv
EOF
cat > /home/learner/solution.sh << 'EOF'
#!/bin/bash
# 使用 case 语句根据扩展名分类文件
EOF
chmod +x /home/learner/solution.sh
