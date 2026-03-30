#!/bin/bash
expected='1. Install dependencies
2. Configure database
3. Run migrations
4. Start development server
5. Run test suite
6. Deploy to staging
7. Run smoke tests
8. Deploy to production'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
