#!/bin/bash
echo "Line 1" > /home/learner/output.txt
echo "Line 2" >> /home/learner/output.txt
echo "Line 3" >> /home/learner/output.txt
echo "Overwritten" > /home/learner/output2.txt
cat /home/learner/output.txt
cat /home/learner/output2.txt
