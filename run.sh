#!/bin/bash

# Function to run the backend
run_backend() {
  echo "Running backend..."
  cd ./backend/cmd/server || exit
  go mod tidy  # Ensure all Go dependencies are downloaded
  go run .
}

# Function to run the frontend
run_frontend() {
  echo "Running frontend..."
  cd ./frontend || exit
  npm install  # Ensure all Node.js dependencies are installed
  npm run dev
}

# Run backend and frontend in parallel
run_backend & run_frontend &

# Wait for both processes to finish
wait