#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/script.sh << 'EOF'
#!/bin/bash
echo "Hello World"
name="Linux"
echo "Welcome to $name"
for i in 1 2 3; do
    echo "Number: $i"
done
exit 0
EOF
