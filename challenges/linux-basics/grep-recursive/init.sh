#!/bin/bash
mkdir -p /home/lab/codebase/{src,tests,docs}
cat > /home/lab/codebase/src/main.py << 'EOF'
def main():
    # TODO: Add error handling
    print("Hello")
    process_data()  # TODO: Optimize this
EOF
cat > /home/lab/codebase/src/utils.py << 'EOF'
def helper():
    return True
EOF
cat > /home/lab/codebase/tests/test_main.py << 'EOF'
def test_main():
    # TODO: Add more test cases
    assert True
EOF
cat > /home/lab/codebase/docs/notes.md << 'EOF'
# Project Notes
- TODO: Update documentation
- Feature complete
EOF
