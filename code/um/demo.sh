#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────────────────────
#  User Management NATS Bridge - Demo Script
#  Shows fire-and-forget message flow:
#    client (nats pub) → UM service → DAL mock (nats sub)
#  A background "nats sub um.>" captures all traffic.
# ─────────────────────────────────────────────────────────────

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

NATS_CLI="${NATS_CLI:-$(go env GOPATH)/bin/nats}"
NATS_SERVER="${NATS_SERVER:-/usr/sbin/nats-server}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TRAFFIC_LOG=$(mktemp /tmp/um-demo-traffic.XXXXXX)
UM_LOG=$(mktemp /tmp/um-demo-service.XXXXXX)
DAL_LOG=$(mktemp /tmp/um-demo-dal.XXXXXX)

# PIDs to clean up
PIDS=()

cleanup() {
    echo ""
    echo -e "${YELLOW}[CLEANUP]${NC} Shutting down..."
    for pid in "${PIDS[@]}"; do
        kill "$pid" 2>/dev/null || true
        wait "$pid" 2>/dev/null || true
    done
    rm -f "$TRAFFIC_LOG" "$UM_LOG" "$DAL_LOG"
    echo -e "${GREEN}[CLEANUP]${NC} Done."
}
trap cleanup EXIT

banner() {
    echo ""
    echo -e "${YELLOW}${BOLD}$1${NC}"
    echo -e "${DIM}───────────────────────────────────────────${NC}"
}

step_banner() {
    echo ""
    echo -e "${CYAN}${BOLD}───────────────────────────────────────────${NC}"
    echo -e "${CYAN}${BOLD}  $1${NC}"
    echo -e "${CYAN}${BOLD}───────────────────────────────────────────${NC}"
}

status_ok() {
    echo -e "  ${GREEN}[OK]${NC} $1"
}

status_fail() {
    echo -e "  ${RED}[FAIL]${NC} $1"
}

wait_for_nats() {
    local retries=20
    while ! "$NATS_CLI" server check connection --server nats://localhost:4222 &>/dev/null; do
        retries=$((retries - 1))
        if [ "$retries" -le 0 ]; then
            status_fail "NATS server did not start"
            exit 1
        fi
        sleep 0.2
    done
}

# ═══════════════════════════════════════════
#  HEADER
# ═══════════════════════════════════════════
echo ""
echo -e "${YELLOW}${BOLD}═══════════════════════════════════════════${NC}"
echo -e "${YELLOW}${BOLD}  User Management NATS Bridge - Demo${NC}"
echo -e "${YELLOW}${BOLD}═══════════════════════════════════════════${NC}"
echo ""
echo -e "${DIM}  Architecture:${NC}"
echo -e "${DIM}  nats pub → UM service → DAL mock (nats sub)${NC}"
echo -e "${DIM}  All traffic captured via: nats sub \"um.>\"${NC}"

# ═══════════════════════════════════════════
#  SETUP
# ═══════════════════════════════════════════
banner "SETUP"

# 1. Start NATS server
echo -ne "  Starting NATS server...                "
"$NATS_SERVER" -p 4222 &>/dev/null &
PIDS+=($!)
wait_for_nats
echo -e "${GREEN}done${NC}"

# 2. Build UM service
echo -ne "  Building UM service...                 "
(cd "$SCRIPT_DIR" && go build -o um . 2>&1) || { status_fail "Build failed"; exit 1; }
echo -e "${GREEN}done${NC}"

# 3. Start UM service
echo -ne "  Starting UM service...                 "
"$SCRIPT_DIR/um" > "$UM_LOG" 2>&1 &
PIDS+=($!)
sleep 1
if ! kill -0 "${PIDS[-1]}" 2>/dev/null; then
    status_fail "UM service crashed on start"
    cat "$UM_LOG"
    exit 1
fi
echo -e "${GREEN}done${NC}"

# 4. Start traffic observer (captures everything on um.>)
echo -ne "  Starting traffic observer...           "
"$NATS_CLI" sub "um.>" > "$TRAFFIC_LOG" 2>&1 &
PIDS+=($!)
sleep 0.3
echo -e "${GREEN}done${NC}"

# 5. Start mock DAL subscribers
echo -ne "  Starting mock DAL responders...        "

dal_subjects=(
    "um.user.db.create"
    "um.user.db.get"
    "um.user.db.update"
    "um.user.db.delete"
    "um.user.db.list"
    "um.role.db.create"
    "um.role.db.get"
    "um.role.db.update"
    "um.role.db.delete"
    "um.role.db.list"
    "um.user.role.db.assign"
    "um.user.role.db.remove"
    "um.user.role.db.list"
)

for subj in "${dal_subjects[@]}"; do
    "$NATS_CLI" sub "$subj" >> "$DAL_LOG" 2>&1 &
    PIDS+=($!)
done
sleep 0.5
echo -e "${GREEN}done${NC}"

echo ""
status_ok "All services running. Starting demo scenarios..."

# ═══════════════════════════════════════════
#  Helper: publish and show
# ═══════════════════════════════════════════
publish_message() {
    local subject="$1"
    local payload="$2"
    local expect_dal_subject="${3:-}"

    echo -e "  ${CYAN}Subject:${NC}  $subject"
    echo -e "  ${CYAN}Payload:${NC}  $payload"
    echo ""

    # Record DAL log size before publish
    local dal_before
    dal_before=$(wc -l < "$DAL_LOG" 2>/dev/null || echo 0)

    # Publish
    "$NATS_CLI" pub "$subject" "$payload" 2>/dev/null

    # Wait briefly for message to flow through UM to DAL
    sleep 0.8

    # Check if DAL received anything new
    local dal_after
    dal_after=$(wc -l < "$DAL_LOG" 2>/dev/null || echo 0)

    if [ "$dal_after" -gt "$dal_before" ]; then
        echo -e "  ${GREEN}=> Message forwarded to DAL${NC}"
        if [ -n "$expect_dal_subject" ]; then
            echo -e "  ${GREEN}=> DAL subject: ${expect_dal_subject}${NC}"
        fi
    else
        echo -e "  ${RED}=> Message NOT forwarded to DAL (validation rejected)${NC}"
    fi
}

# ═══════════════════════════════════════════
#  DEMO SCENARIOS
# ═══════════════════════════════════════════

# --- 1. Create User ---
step_banner "1. Create User"
publish_message "um.user.create" \
    '{"user":{"name":"John","email":"john@motadata.com"},"timestamp":"2026-03-10T12:00:00Z"}' \
    "um.user.db.create"

sleep 1

# --- 2. Get User ---
step_banner "2. Get User"
publish_message "um.user.get" \
    '{"user_id":"abc-123","timestamp":"2026-03-10T12:00:01Z"}' \
    "um.user.db.get"

sleep 1

# --- 3. List Users ---
step_banner "3. List Users"
publish_message "um.user.list" \
    '{"limit":10,"offset":0,"timestamp":"2026-03-10T12:00:02Z"}' \
    "um.user.db.list"

sleep 1

# --- 4. Update User ---
step_banner "4. Update User"
publish_message "um.user.update" \
    '{"user":{"id":"abc-123","name":"John Updated","email":"john.updated@motadata.com"},"timestamp":"2026-03-10T12:00:03Z"}' \
    "um.user.db.update"

sleep 1

# --- 5. Create Role ---
step_banner "5. Create Role"
publish_message "um.role.create" \
    '{"role":{"name":"Admin","permissions":["user:create","role:assign"]},"timestamp":"2026-03-10T12:00:04Z"}' \
    "um.role.db.create"

sleep 1

# --- 6. Assign Role to User ---
step_banner "6. Assign Role to User"
publish_message "um.user.role.assign" \
    '{"user_id":"abc-123","role_id":"role-456","timestamp":"2026-03-10T12:00:05Z"}' \
    "um.user.role.db.assign"

sleep 1

# --- 7. List User Roles ---
step_banner "7. List User Roles"
publish_message "um.user.role.list" \
    '{"user_id":"abc-123","limit":10,"offset":0,"timestamp":"2026-03-10T12:00:06Z"}' \
    "um.user.role.db.list"

sleep 1

# --- 8. Validation Error: Missing Fields ---
step_banner "8. Validation Error - Missing User Fields"
echo -e "  ${DIM}(UM service should reject this - name and email are empty)${NC}"
echo ""
publish_message "um.user.create" \
    '{"user":{"name":"","email":""},"timestamp":"2026-03-10T12:00:07Z"}'

sleep 1

# --- 9. Validation Error: Missing User ID ---
step_banner "9. Validation Error - Missing User ID"
echo -e "  ${DIM}(UM service should reject this - user_id is empty)${NC}"
echo ""
publish_message "um.user.get" \
    '{"user_id":"","timestamp":"2026-03-10T12:00:08Z"}'

sleep 1

# --- 10. Delete User ---
step_banner "10. Delete User"
publish_message "um.user.delete" \
    '{"user_id":"abc-123","timestamp":"2026-03-10T12:00:09Z"}' \
    "um.user.db.delete"

sleep 1

# --- 11. Remove Role ---
step_banner "11. Remove Role from User"
publish_message "um.user.role.remove" \
    '{"user_id":"abc-123","role_id":"role-456","timestamp":"2026-03-10T12:00:10Z"}' \
    "um.user.role.db.remove"

sleep 1

# ═══════════════════════════════════════════
#  TRAFFIC SUMMARY
# ═══════════════════════════════════════════
banner "TRAFFIC SUMMARY"

echo -e "  ${CYAN}All messages observed on um.> :${NC}"
echo ""
if [ -s "$TRAFFIC_LOG" ]; then
    # Show subject lines from nats sub output
    grep -E "^\[#[0-9]+\]" "$TRAFFIC_LOG" | while IFS= read -r line; do
        echo -e "  ${DIM}${line}${NC}"
    done
else
    echo -e "  ${DIM}(no traffic captured)${NC}"
fi

echo ""
echo -e "  ${CYAN}Messages received by mock DAL:${NC}"
echo ""
if [ -s "$DAL_LOG" ]; then
    grep -E "^\[#[0-9]+\]" "$DAL_LOG" | while IFS= read -r line; do
        echo -e "  ${DIM}${line}${NC}"
    done
else
    echo -e "  ${DIM}(no DAL messages)${NC}"
fi

echo ""

# ═══════════════════════════════════════════
#  UM SERVICE LOGS
# ═══════════════════════════════════════════
banner "UM SERVICE LOGS"

if [ -s "$UM_LOG" ]; then
    while IFS= read -r line; do
        echo -e "  ${DIM}${line}${NC}"
    done < "$UM_LOG"
else
    echo -e "  ${DIM}(no logs)${NC}"
fi

echo ""
echo -e "${GREEN}${BOLD}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}${BOLD}  Demo complete!${NC}"
echo -e "${GREEN}${BOLD}═══════════════════════════════════════════${NC}"
echo ""
