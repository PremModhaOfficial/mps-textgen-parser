#!/usr/bin/env bash
set -euo pipefail

# ─── Config ───
MPS_OUTPUT="$HOME/UserManagmentMps/languages/UserManagement.sandbox/source_gen/src"
SDK_SRC="$HOME/projects/nextgen/gosdk/motadata-go-sdk/src/motadatagosdk"
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$PROJECT_DIR/src"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# ─── Sync generated files (.go + .sql) into src/ ───
echo -ne "${CYAN}Syncing generated files...${NC} "
mkdir -p "$SRC_DIR"
count=0
for f in "$MPS_OUTPUT"/*.go "$MPS_OUTPUT"/*.sql; do
    [ -f "$f" ] || continue
    cp "$f" "$SRC_DIR/"
    count=$((count + 1))
done
echo -e "${GREEN}${count} files copied${NC}"

# ─── Extract hooks from generated files and merge into userDefinedHooks.go ───
echo -e "${CYAN}Merging hooks...${NC}"
python3 "$PROJECT_DIR/merge_hooks.py"

# ─── Extract project name from main.go first line comment ───
# Expects: // <projectName>
PROJECT_NAME=$(head -1 "$SRC_DIR/main.go" | sed 's|^// *||' | tr '[:upper:]' '[:lower:]')
if [ -z "$PROJECT_NAME" ]; then
    echo -e "${RED}Could not extract project name from main.go${NC}"
    exit 1
fi
echo -e "${CYAN}Project:${NC} ${PROJECT_NAME}"

# ─── Sync SDK ───
echo -ne "${CYAN}Syncing SDK...${NC} "
if [ -d "$SDK_SRC" ]; then
    rsync -a --delete "$SDK_SRC/" "$PROJECT_DIR/motadata-go-sdk/"
    # Patch module path: source uses "motadatagosdk" but our go.mod expects "dev.azure.com/..."
    sed -i 's|^module motadatagosdk|module dev.azure.com/Motadata/NextGen/motadata-go-sdk|' \
        "$PROJECT_DIR/motadata-go-sdk/go.mod"
    # Fix SDK internal imports to match the module path
    find "$PROJECT_DIR/motadata-go-sdk" -name '*.go' -exec \
        sed -i 's|"motadatagosdk/|"dev.azure.com/Motadata/NextGen/motadata-go-sdk/|g' {} +
    echo -e "${GREEN}ok (patched module path)${NC}"
else
    echo -e "${RED}SDK not found at $SDK_SRC${NC}"
    exit 1
fi

# ─── Summary ───
echo ""
echo -e "${CYAN}Project ready at:${NC} $PROJECT_DIR"
echo -e "${CYAN}src/:${NC}"
ls -1 "$SRC_DIR" | while read -r f; do
    echo "  $f"
done
echo ""
echo -e "${GREEN}Starting ${PROJECT_NAME}...${NC}"

# Free ports before starting
fuser -k 4229/tcp 8222/tcp 2>/dev/null || true

docker compose build --build-arg "APP_NAME=$PROJECT_NAME" && \
    docker compose up
