#!/bin/bash
out1=$(bash /home/learner/solution.sh 1 2>/dev/null)
out2=$(bash /home/learner/solution.sh 2 2>/dev/null)
out3=$(bash /home/learner/solution.sh 3 2>/dev/null)
out4=$(bash /home/learner/solution.sh 4 2>/dev/null)
out5=$(bash /home/learner/solution.sh 9 2>/dev/null)
if [ "$out1" = "You selected: View date" ] && \
   [ "$out2" = "You selected: View user" ] && \
   [ "$out3" = "You selected: View directory" ] && \
   [ "$out4" = "Goodbye!" ] && \
   [ "$out5" = "Invalid option" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "Option 1: '$out1'"
    echo "Option 2: '$out2'"
    echo "Option 3: '$out3'"
    echo "Option 4: '$out4'"
    echo "Option 9: '$out5'"
    exit 1
fi
