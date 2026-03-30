#!/bin/bash
mkdir -p /home/lab
cat > /home/lab/emails.txt << 'EOF'
Contact us at support@example.com for help.
This line has no email.
Send reports to admin@company.org by Friday.
Phone: 555-1234
The developer john.doe@dev-team.io fixed the bug.
Address: 123 Main Street
CC: marketing+newsletter@bigcorp.com
No contact info here.
EOF
