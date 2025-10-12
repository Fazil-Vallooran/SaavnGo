#!/bin/bash

# JioSaavn API Runner Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default port
PORT=${SERVER_PORT:-8080}

echo -e "${GREEN}Starting JioSaavn API...${NC}"

# Check if port is in use
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo -e "${RED}Error: Port $PORT is already in use.${NC}"
    echo -e "${YELLOW}Try running with a different port:${NC}"
    echo -e "  SERVER_PORT=3000 ./run.sh"
    exit 1
fi

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Install dependencies if needed
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    echo -e "${YELLOW}Installing dependencies...${NC}"
    go mod download
fi

# Run the application
echo -e "${GREEN}Starting server on port $PORT...${NC}"
echo -e "${GREEN}Health check: http://localhost:$PORT/health${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""

go run main.go
