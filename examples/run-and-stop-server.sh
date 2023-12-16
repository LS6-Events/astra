#!/bin/bash

# Start the Go web server in the background on port 8000
go run . &
SERVER_PID=$!

# Wait for the server to start
# This is a simple loop that checks if the port 8000 is being used
while ! lsof -i:8000 -sTCP:LISTEN -t >/dev/null; do
    echo "Waiting for server to start..."
    sleep 1
done

echo "Server started."

# Optional: Wait for a specific time before killing the server
# sleep 10

# Kill the server process
kill "$(lsof -ti:8000)"
echo "Server stopped."