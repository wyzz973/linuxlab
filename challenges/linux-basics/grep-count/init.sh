#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/survey.txt << 'EOF'
User1: The app works great!
User2: I got an Error when clicking submit
User3: Everything is fine
User4: error loading the page
User5: Great product, no issues
User6: ERROR: cannot connect to server
User7: Minor error in the calculation
User8: Love this app!
User9: The ERROR message is confusing
User10: All good here
EOF
