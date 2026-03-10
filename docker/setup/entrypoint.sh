#!/bin/bash
set -e

# Color codes for output
WARN='\033[1;33m'
OK='\033[0;32m'
NC='\033[0m' # No Color

# Flag to track if any check failed
ANY_FAILED=false

echo "Starting Filestash setup checks..."
echo ""

# 1. Check filestash user owns /app/data
echo -n "Checking /app/data ownership... "
if [ "$(stat -c '%U' /app/data)" = "filestash" ]; then
    echo -e "${OK}OK${NC}"
else
    echo -e "${WARN}WARNING${NC}"
    echo "  User filestash does not own /app/data. Current owner: $(stat -c '%U' /app/data)"
    ANY_FAILED=true
fi

# 2. Check required environment variables
echo -n "Checking FILESTASH_S3_ENDPOINT... "
if [ -z "$FILESTASH_S3_ENDPOINT" ]; then
    echo -e "${WARN}WARNING${NC} (not set)"
    ANY_FAILED=true
else
    echo -e "${OK}OK${NC}"
fi

echo -n "Checking FILESTASH_S3_ACCESS_KEY... "
if [ -z "$FILESTASH_S3_ACCESS_KEY" ]; then
    echo -e "${WARN}WARNING${NC} (not set)"
    ANY_FAILED=true
else
    echo -e "${OK}OK${NC}"
fi

echo -n "Checking FILESTASH_S3_SECRET_ACCESS_KEY... "
if [ -z "$FILESTASH_S3_SECRET_ACCESS_KEY" ]; then
    echo -e "${WARN}WARNING${NC} (not set)"
    ANY_FAILED=true
else
    echo -e "${OK}OK${NC}"
fi

# 3. Check and setup config file
CONFIG_DIR="/app/data/state/config"
CONFIG_FILE="${CONFIG_DIR}/config.json"
echo -n "Checking config file... "
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${WARN}NOT FOUND${NC}"
    echo "  Creating config directory and copying default config..."
    mkdir -p "$CONFIG_DIR"
    if [ -f "/app/setup/config.json" ]; then
        cp "/app/setup/config.json" "$CONFIG_FILE"
        chown filestash:filestash "$CONFIG_FILE"
        echo "  Default config copied to $CONFIG_FILE"
    else
        echo -e "${WARN}WARNING${NC}: /app/setup/config.json not found"
        ANY_FAILED=true
    fi
else
    echo -e "${OK}OK${NC}"
fi

# 4. Check ca-certificates is installed
echo -n "Checking ca-certificates... "
if dpkg -l | grep -q '^ii  ca-certificates'; then
    echo -e "${OK}OK${NC}"
else
    echo -e "${WARN}WARNING${NC} (not installed)"
    ANY_FAILED=true
fi

echo ""

# Summary
if [ "$ANY_FAILED" = false ]; then
    echo -e "${OK}All checks passed!${NC}"
else
    echo -e "${WARN}Some checks failed or raised warnings. Check output above.${NC}"
fi

echo ""
echo "Starting Filestash..."
exec /app/filestash "$@"
