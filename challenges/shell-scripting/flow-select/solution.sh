#!/bin/bash
case "$1" in
    1) echo "You selected: View date" ;;
    2) echo "You selected: View user" ;;
    3) echo "You selected: View directory" ;;
    4) echo "Goodbye!" ;;
    *) echo "Invalid option" ;;
esac
