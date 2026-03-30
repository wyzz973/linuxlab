#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/code.py << 'PYEOF'
# This is a Python script
# Author: admin

def calculate(x, y):
    # Add two numbers
    return x + y

def greet(name):
    # Print greeting
    print(f"Hello {name}")

# Main entry point
if __name__ == "__main__":
    result = calculate(3, 4)
    greet("World")
PYEOF
