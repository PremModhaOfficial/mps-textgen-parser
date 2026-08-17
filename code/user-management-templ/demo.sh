#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────────────────────
#  MPS-Generated NATS Bridge - Demo Script
#  Tests the DSL-generated code (motadata.* subjects)
#  Includes password validation test via types.MeetsPolicy
# ─────────────────────────────────────────────────────────────

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

NATS_CLI="${NATS_CLI:-$(go env GOPATH)/bin/nats}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
S="nats://localhost:4229"
TRAFFIC_LOG=$(mktemp /tmp/gen-demo-traffic.XXXXXX)
DAL_LOG=$(mktemp /tmp/gen-demo-dal.XXXXXX)

PIDS=()

cleanup() {
    echo ""
    echo -e "${YELLOW}[CLEANUP]${NC} Shutting down..."
    for pid in "${PIDS[@]}"; do
        kill "$pid" 2>/dev/null || true
        wait "$pid" 2>/dev/null || true
    done
    rm -f "$TRAFFIC_LOG" "$DAL_LOG"
    echo -e "${GREEN}[CLEANUP]${NC} Done."
}
trap cleanup EXIT

step() {
    echo ""
    echo -e "${CYAN}${BOLD}─── $1 ───${NC}"
}

wait_for_nats() {
    local retries=20
    while ! "$NATS_CLI" server check connection --server "$S" &>/dev/null; do
        retries=$((retries - 1))
        if [ "$retries" -le 0 ]; then echo -e "${RED}NATS not reachable at $S${NC}"; exit 1; fi
        sleep 0.2
    done
}

publish_message() {
    local subject="$1" payload="$2" expect_dal="${3:-}"
    echo -e "  ${CYAN}Subject:${NC} $subject"
    echo -e "  ${CYAN}Payload:${NC} $payload"
    local reply
    if reply=$("$NATS_CLI" req "$subject" "$payload" --server "$S" --timeout 5s 2>&1); then
        echo -e "  ${GREEN}=> Got reply${NC}${expect_dal:+ ($expect_dal)}"
        echo -e "  ${DIM}${reply}${NC}"
    else
        echo -e "  ${RED}=> Error/Rejected${NC}"
        echo -e "  ${DIM}${reply}${NC}"
    fi
}

# ═══════════════════════════════════════════
echo -e "${YELLOW}${BOLD}═══ MPS-Generated NATS Bridge Demo ═══${NC}"
echo -e "${DIM}  Service: docker compose on ${S}${NC}"
echo -e "${DIM}  Subjects: motadata.*${NC}"

# ─── SETUP ───
step "SETUP"

echo -ne "  Checking NATS...     "
wait_for_nats
echo -e "${GREEN}connected to ${S}${NC}"

echo -ne "  Traffic observer...  "
"$NATS_CLI" sub "motadata.>" --server "$S" > "$TRAFFIC_LOG" 2>&1 &
PIDS+=($!)
sleep 0.3
echo -e "${GREEN}ok${NC}"

echo -ne "  Mock DAL...          "
"$NATS_CLI" reply "motadata.user.db.>" '{"status":"ok"}' --server "$S" >> "$DAL_LOG" 2>&1 &
PIDS+=($!)
"$NATS_CLI" reply "motadata.roles.db.>" '{"status":"ok"}' --server "$S" >> "$DAL_LOG" 2>&1 &
PIDS+=($!)
"$NATS_CLI" reply "motadata.user.roles.db.>" '{"status":"ok"}' --server "$S" >> "$DAL_LOG" 2>&1 &
PIDS+=($!)
"$NATS_CLI" reply "motadata.permissions.db.>" '{"status":"ok"}' --server "$S" >> "$DAL_LOG" 2>&1 &
PIDS+=($!)
"$NATS_CLI" reply "motadata.roles.permissions.db.>" '{"status":"ok"}' --server "$S" >> "$DAL_LOG" 2>&1 &
PIDS+=($!)
sleep 0.5
echo -e "${GREEN}ok (5 mock DAL endpoints)${NC}"

# ─── TESTS ───

step "1. Create User (valid + strong password)"
publish_message "motadata.user.create" \
    '{"user":{"user_name":"John","password":"Str0ng!Pass","user_mail":"john@example.com"},"timestamp":"2026-03-10T12:00:00Z"}' \
    "motadata.user.db.create"

step "2. Create User (weak password - should REJECT)"
publish_message "motadata.user.create" \
    '{"user":{"user_name":"Jane","password":"weak"},"timestamp":"2026-03-10T12:00:01Z"}'

step "3. Create User (missing fields - should REJECT)"
publish_message "motadata.user.create" \
    '{"user":{"user_name":""},"timestamp":"2026-03-10T12:00:02Z"}'

step "4. Get User"
publish_message "motadata.user.get" \
    '{"user_id":"abc-123","timestamp":"2026-03-10T12:00:03Z"}' \
    "motadata.user.db.get"

step "5. List Users"
publish_message "motadata.user.list" \
    '{"limit":10,"offset":0,"timestamp":"2026-03-10T12:00:04Z"}' \
    "motadata.user.db.list"

step "6. Update User"
publish_message "motadata.user.update" \
    '{"user":{"id":"abc-123","user_name":"John Updated"},"timestamp":"2026-03-10T12:00:05Z"}' \
    "motadata.user.db.update"

step "7. Create Role"
publish_message "motadata.roles.create" \
    '{"roles":{"name":"Admin","given_to":"devs"},"timestamp":"2026-03-10T12:00:06Z"}' \
    "motadata.roles.db.create"

step "8. Assign Role to User"
publish_message "motadata.user.roles.assign" \
    '{"user_id":"abc-123","roles_id":"role-456","timestamp":"2026-03-10T12:00:07Z"}' \
    "motadata.user.roles.db.assign"

step "9. List User Roles"
publish_message "motadata.user.roles.list" \
    '{"user_id":"abc-123","limit":10,"offset":0,"timestamp":"2026-03-10T12:00:08Z"}' \
    "motadata.user.roles.db.list"

step "10. Remove Role from User"
publish_message "motadata.user.roles.remove" \
    '{"user_id":"abc-123","roles_id":"role-456","timestamp":"2026-03-10T12:00:09Z"}' \
    "motadata.user.roles.db.remove"

step "11. Delete User"
publish_message "motadata.user.delete" \
    '{"user_id":"abc-123","timestamp":"2026-03-10T12:00:10Z"}' \
    "motadata.user.db.delete"

# ─── Permissions CRUD ───

step "12. Create Permission"
publish_message "motadata.permissions.create" \
    '{"permissions":{"id":"perm-1","name":"view_dashboard","moduele":"dashboard","can":"read"},"timestamp":"2026-03-10T12:00:11Z"}' \
    "motadata.permissions.db.create"

step "13. Create Permission (missing fields - should REJECT)"
publish_message "motadata.permissions.create" \
    '{"permissions":{"id":"perm-2","name":""},"timestamp":"2026-03-10T12:00:12Z"}'

step "14. Get Permission"
publish_message "motadata.permissions.get" \
    '{"permissions_id":"perm-1","timestamp":"2026-03-10T12:00:13Z"}' \
    "motadata.permissions.db.get"

step "15. List Permissions"
publish_message "motadata.permissions.list" \
    '{"limit":10,"offset":0,"timestamp":"2026-03-10T12:00:14Z"}' \
    "motadata.permissions.db.list"

step "16. Update Permission"
publish_message "motadata.permissions.update" \
    '{"permissions":{"id":"perm-1","name":"edit_dashboard","moduele":"dashboard","can":"write"},"timestamp":"2026-03-10T12:00:15Z"}' \
    "motadata.permissions.db.update"

step "17. Delete Permission"
publish_message "motadata.permissions.delete" \
    '{"permissions_id":"perm-1","timestamp":"2026-03-10T12:00:16Z"}' \
    "motadata.permissions.db.delete"

# ─── Roles <-> Permissions Relation ───

step "18. Assign Permission to Role"
publish_message "motadata.roles.permissions.assign" \
    '{"roles_id":"role-456","permissions_id":"perm-1","timestamp":"2026-03-10T12:00:17Z"}' \
    "motadata.roles.permissions.db.assign"

step "19. List Role Permissions"
publish_message "motadata.roles.permissions.list" \
    '{"roles_id":"role-456","limit":10,"offset":0,"timestamp":"2026-03-10T12:00:18Z"}' \
    "motadata.roles.permissions.db.list"

step "20. Remove Permission from Role"
publish_message "motadata.roles.permissions.remove" \
    '{"roles_id":"role-456","permissions_id":"perm-1","timestamp":"2026-03-10T12:00:19Z"}' \
    "motadata.roles.permissions.db.remove"

step "21. Service Discovery"
echo -e "  ${CYAN}Subject:${NC} \$SRV.INFO"
reply=$("$NATS_CLI" req '$SRV.INFO' '' --server "$S" --timeout 2s 2>&1) || true
echo -e "  ${DIM}${reply}${NC}"


# ─── SUMMARY ───
step "TRAFFIC SUMMARY"
echo -e "  ${CYAN}All messages on motadata.> :${NC}"
grep -E "^\[#[0-9]+\]" "$TRAFFIC_LOG" 2>/dev/null | while IFS= read -r line; do
    echo -e "  ${DIM}${line}${NC}"
done

echo ""
echo -e "  ${CYAN}DAL received:${NC}"
grep -E "^\[#[0-9]+\]" "$DAL_LOG" 2>/dev/null | while IFS= read -r line; do
    echo -e "  ${DIM}${line}${NC}"
done

step "SERVICE LOGS"
while IFS= read -r line; do echo -e "  ${DIM}${line}${NC}"; done < "$UM_LOG"

echo ""
echo -e "${GREEN}${BOLD}═══ Demo complete! ═══${NC}"
