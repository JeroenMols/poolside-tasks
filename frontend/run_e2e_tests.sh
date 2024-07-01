#!/bin/bash

# Start the backend server
(cd ../backend && ./start_server.sh)&
BACKEND_PID=$!

# Clear the app cached of the electron app
rm -r ~/Library/Application\ Support/frontend-electron

# Build and run tests
npm run build
npx playwright test --reporter=html

# Kill the backend server
kill $BACKEND_PID &> /dev/null
kill $(pgrep -f go-build) > /dev/null
    