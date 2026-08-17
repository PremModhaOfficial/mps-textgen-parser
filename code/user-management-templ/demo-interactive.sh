#!/usr/bin/env bash
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
NATS_CLI="${NATS_CLI:-$(go env GOPATH)/bin/nats}"
S="nats://localhost:4229"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BLUE='\033[0;34m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

PIDS=()
PASS=0
FAIL=0
TOTAL_RTT=0
LOG_FILE=$(mktemp)
LOG_POS=0
LOG_TAIL_PID=""

cleanup() {
    [ -n "$LOG_TAIL_PID" ] && kill "$LOG_TAIL_PID" 2>/dev/null || true
    for pid in "${PIDS[@]}"; do kill "$pid" 2>/dev/null || true; done
    rm -f "$LOG_FILE"
}
trap cleanup EXIT

show_logs() {
    sleep 0.5  # give async hooks time to log
    local current_size
    current_size=$(wc -c < "$LOG_FILE" 2>/dev/null) || current_size=0
    if [ "$current_size" -gt "$LOG_POS" ]; then
        local new_lines
        new_lines=$(tail -c +"$((LOG_POS + 1))" "$LOG_FILE" 2>/dev/null | grep -v '^$' | grep -v '^app ' | tail -12) || true
        if [ -n "$new_lines" ]; then
            echo ""
            echo -e "  ${BLUE}${DIM}┄┄ app logs ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄${NC}"
            while IFS= read -r line; do
                echo -e "  ${BLUE}${DIM}[LOG]${NC} ${DIM}${line}${NC}"
            done <<< "$new_lines"
            echo -e "  ${BLUE}${DIM}┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄${NC}"
        fi
        LOG_POS=$current_size
    fi
}

pause() {
    echo ""
    echo -ne "${DIM}Press Enter to continue...${NC}"
    read -r
    echo ""
}

send() {
    local label="$1" subject="$2" payload="$3" expect="${4:-}" hooks="${5:-}"
    echo -e "${CYAN}${BOLD}━━━ $label ━━━${NC}"
    echo ""
    echo -e "  ${YELLOW}Subject:${NC}  $subject"
    echo -e "  ${YELLOW}Payload:${NC}"
    echo "$payload" | python3 -m json.tool 2>/dev/null | while IFS= read -r line; do
        echo -e "    ${DIM}$line${NC}"
    done || true
    echo ""
    [ -n "$expect" ] && echo -e "  ${DIM}Expect: $expect${NC}"
    [ -n "$hooks" ] && echo -e "  ${MAGENTA}Hooks:${NC}  $hooks"
    echo ""

    local reply rtt_line rtt_ms
    if reply=$("$NATS_CLI" req "$subject" "$payload" --server "$S" --timeout 5s 2>&1); then
        rtt_line=$(echo "$reply" | grep -oP 'rtt \K[0-9.]+m?s' | head -1 || true)
        if echo "$reply" | grep -q "Nats-Service-Error" 2>/dev/null; then
            local err=$(echo "$reply" | grep "Nats-Service-Error:" | head -1 | sed 's/.*Error: //' || true)
            local code=$(echo "$reply" | grep "Error-Code:" | head -1 | sed 's/.*Code: //' || true)
            echo -e "  ${RED}${BOLD}REJECTED ${code}${NC} ${err}"
            if [[ "$expect" == *"reject"* ]] || [[ "$expect" == *"400"* ]]; then
                echo -e "  ${GREEN}(expected rejection)${NC}"
                PASS=$((PASS + 1))
            else
                FAIL=$((FAIL + 1))
            fi
        else
            echo -e "  ${GREEN}${BOLD}OK${NC} → forwarded to DAL  ${DIM}(${rtt_line:-?})${NC}"
            local data
            data=$(echo "$reply" | tail -1 || true)
            [ -n "$data" ] && echo -e "  ${DIM}Reply: $data${NC}"
            PASS=$((PASS + 1))
        fi
    else
        echo -e "  ${RED}ERROR${NC}"
        echo -e "  ${DIM}$reply${NC}"
        FAIL=$((FAIL + 1))
    fi
}

stress() {
    local label="$1" subject="$2" payload="$3" count="${4:-50}"
    echo -e "${MAGENTA}${BOLD}━━━ STRESS: $label ($count requests) ━━━${NC}"
    echo ""

    local start_time end_time elapsed ok=0 err=0
    start_time=$(date +%s%N)

    for ((i=1; i<=count; i++)); do
        if "$NATS_CLI" req "$subject" "$payload" --server "$S" --timeout 5s &>/dev/null; then
            ok=$((ok + 1))
        else
            err=$((err + 1))
        fi
    done

    end_time=$(date +%s%N)
    elapsed=$(( (end_time - start_time) / 1000000 ))
    local rps=0
    if [ "$elapsed" -gt 0 ]; then
        rps=$(( count * 1000 / elapsed ))
    fi

    echo -e "  ${GREEN}${BOLD}${ok}${NC} OK  ${RED}${BOLD}${err}${NC} ERR  ${DIM}in ${elapsed}ms (~${rps} req/s)${NC}"
    echo ""
}

# ═══════════════════════════════════════════
clear
echo -e "${YELLOW}${BOLD}"
echo "  ╔══════════════════════════════════════════════════════╗"
echo "  ║   MPS-Generated NATS Microservice Demo              ║"
echo "  ║   DSL → Go Code → Docker → Live Service             ║"
echo "  ║                                                      ║"
echo "  ║   Features: Async/Sync Hooks + OTEL Debug Logging   ║"
echo "  ╚══════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo -e "  ${DIM}Service running on ${S}${NC}"
echo -e "  ${DIM}Subjects: motadata.*${NC}"
echo -e "  ${DIM}Tip: run 'docker compose logs -f app' in another terminal for live hook logs${NC}"
echo ""

# Check connection
echo -ne "  ${CYAN}Checking NATS connection...${NC} "
if "$NATS_CLI" server check connection --server "$S" &>/dev/null; then
    echo -e "${GREEN}connected${NC}"
else
    echo -e "${RED}cannot connect to $S${NC}"
    exit 1
fi

# Start mock DAL responders with realistic data
echo -ne "  ${CYAN}Starting mock DAL responders...${NC} "

# User mocks (--count=0 = unlimited replies)
"$NATS_CLI" reply "motadata.user.db.create" '{"status":"created","user":{"id":"abc-123","user_name":"John","user_mail":"john@example.com","created_at":"2026-03-12T12:00:01Z"}}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)
"$NATS_CLI" reply "motadata.user.db.update" '{"status":"updated","user":{"id":"abc-123","user_name":"John Updated","user_mail":"john.new@example.com"}}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)
"$NATS_CLI" reply "motadata.user.db.get" '{"user":{"id":"abc-123","user_name":"John","user_mail":"john@example.com","created_at":"2026-03-12T12:00:01Z"}}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)
"$NATS_CLI" reply "motadata.user.db.list" '{"users":[{"id":"abc-123","user_name":"John","user_mail":"john@example.com"},{"id":"def-456","user_name":"Jane","user_mail":"jane@example.com"}],"total":2}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)
"$NATS_CLI" reply "motadata.user.db.delete" '{"status":"deleted","user_id":"abc-123"}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)

# Roles mocks
"$NATS_CLI" reply "motadata.roles.db.create" '{"status":"created","roles":{"id":"role-456","name":"Admin","desc":"Full access role"}}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)
"$NATS_CLI" reply "motadata.roles.db.update" '{"status":"updated","roles":{"id":"role-456","name":"SuperAdmin","desc":"Full access role with audit"}}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)

# UserRoles relation mocks
"$NATS_CLI" reply "motadata.user.roles.db.assign" '{"status":"assigned","user_id":"abc-123","roles_id":"role-456"}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)
"$NATS_CLI" reply "motadata.user.roles.db.list" '{"roles":[{"id":"role-456","name":"Admin","desc":"Full access role"}],"total":1}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)
"$NATS_CLI" reply "motadata.user.roles.db.remove" '{"status":"removed","user_id":"abc-123","roles_id":"role-456"}' --server "$S" --count=0 &>/dev/null &
PIDS+=($!)

sleep 0.5
echo -e "${GREEN}ready (10 mock endpoints, unlimited replies)${NC}"

# Start tailing docker app logs in background
docker compose -f "$SCRIPT_DIR/docker-compose.yml" logs -f --no-color app > "$LOG_FILE" 2>&1 &
LOG_TAIL_PID=$!
sleep 0.3  # give it a moment to open the log stream

pause

# ─── PHASE 1: FUNCTIONAL TESTS ───
echo -e "${YELLOW}${BOLD}  ── Phase 1: Functional Tests (11 endpoints) ──${NC}"
echo ""

send "1/11  Create Role" \
    "motadata.roles.create" \
    '{"roles":{"name":"Admin","desc":"Full access role"},"timestamp":"2026-03-12T12:00:00Z"}' \
    "→ motadata.roles.db.create" \
    "none"
show_logs

pause

send "2/11  Update Role" \
    "motadata.roles.update" \
    '{"roles":{"id":"role-456","name":"SuperAdmin","desc":"Full access role with audit"},"timestamp":"2026-03-12T12:00:01Z"}' \
    "→ motadata.roles.db.update" \
    "none"
show_logs

pause

send "3/11  Create User (valid)" \
    "motadata.user.create" \
    '{"user":{"user_name":"John","password":"Str0ngPass1","user_mail":"john@example.com"},"timestamp":"2026-03-12T12:00:02Z"}' \
    "→ motadata.user.db.create" \
    "pre: D(sync) C(async) B(sync) A(async) | post: notifyAdmin(async) sendWelcome(async)"
show_logs

pause

send "4/11  Create User (missing fields → should REJECT)" \
    "motadata.user.create" \
    '{"user":{"user_name":""},"timestamp":"2026-03-12T12:00:03Z"}' \
    "400 rejection (validation)" \
    "pre-hooks run first, then validation rejects"
show_logs

pause

send "5/11  Update User" \
    "motadata.user.update" \
    '{"user":{"id":"abc-123","user_name":"John Updated","user_mail":"john.new@example.com"},"timestamp":"2026-03-12T12:00:04Z"}' \
    "→ motadata.user.db.update" \
    "pre: notifyCache(async) auditLog(async) | post: invalidate(async)"
show_logs

pause

send "6/11  Get User by ID" \
    "motadata.user.get" \
    '{"user_id":"abc-123","timestamp":"2026-03-12T12:00:05Z"}' \
    "→ motadata.user.db.get" \
    "pre: traceAcces(async) | post: filterPII(sync) cacheResult(async) enrichData(sync)"
show_logs

pause

send "7/11  List All Users" \
    "motadata.user.list" \
    '{"limit":10,"offset":0,"timestamp":"2026-03-12T12:00:06Z"}' \
    "→ motadata.user.db.list" \
    "none (no hooks)"
show_logs

pause

send "8/11  Assign Role to User" \
    "motadata.user.roles.assign" \
    '{"user_id":"abc-123","roles_id":"role-456","timestamp":"2026-03-12T12:00:07Z"}' \
    "→ motadata.user.roles.db.assign" \
    "pre/post relation hooks"
show_logs

pause

send "9/11  List User Roles" \
    "motadata.user.roles.list" \
    '{"user_id":"abc-123","limit":10,"offset":0,"timestamp":"2026-03-12T12:00:08Z"}' \
    "→ motadata.user.roles.db.list" \
    "pre/post relation hooks"
show_logs

pause

send "10/11  Remove Role from User" \
    "motadata.user.roles.remove" \
    '{"user_id":"abc-123","roles_id":"role-456","timestamp":"2026-03-12T12:00:09Z"}' \
    "→ motadata.user.roles.db.remove" \
    "pre/post relation hooks"
show_logs

pause

send "11/11  Delete User" \
    "motadata.user.delete" \
    '{"user_id":"abc-123","timestamp":"2026-03-12T12:00:10Z"}' \
    "→ motadata.user.db.delete" \
    "pre: Xyz(sync) Abc(sync) validateTen(sync) | post: mustDelete(sync)"
show_logs

echo ""
echo -e "${CYAN}${BOLD}  ── Phase 1 Results: ${GREEN}${PASS} passed${NC} ${RED}${FAIL} failed${NC} ──"
echo ""

pause

# ─── PHASE 2: STRESS TESTS ───
echo -e "${YELLOW}${BOLD}  ── Phase 2: Stress Tests ──${NC}"
echo -e "  ${DIM}Sequential requests to measure throughput & stability${NC}"
echo ""

stress "Create User (mixed sync+async hooks)" \
    "motadata.user.create" \
    '{"user":{"user_name":"StressUser","password":"Pass123","user_mail":"stress@test.com"},"timestamp":"2026-03-16T00:00:00Z"}' \
    100

stress "Get User (async pre + mixed post hooks)" \
    "motadata.user.get" \
    '{"user_id":"abc-123","timestamp":"2026-03-16T00:00:01Z"}' \
    100

stress "Delete User (all sync hooks)" \
    "motadata.user.delete" \
    '{"user_id":"abc-123","timestamp":"2026-03-16T00:00:02Z"}' \
    100

stress "List Users (no hooks)" \
    "motadata.user.list" \
    '{"limit":10,"offset":0,"timestamp":"2026-03-16T00:00:03Z"}' \
    100

stress "Update User (all async hooks)" \
    "motadata.user.update" \
    '{"user":{"id":"abc-123","user_name":"Updated","user_mail":"u@test.com"},"timestamp":"2026-03-16T00:00:04Z"}' \
    100

stress "Assign Role (relation)" \
    "motadata.user.roles.assign" \
    '{"user_id":"abc-123","roles_id":"role-456","timestamp":"2026-03-16T00:00:05Z"}' \
    100

echo ""
echo -e "${GREEN}${BOLD}  ╔══════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}${BOLD}  ║              Demo + Stress Test Complete!            ║${NC}"
echo -e "${GREEN}${BOLD}  ║                                                      ║${NC}"
echo -e "${GREEN}${BOLD}  ║  Functional: ${PASS} passed / ${FAIL} failed                    ║${NC}"
echo -e "${GREEN}${BOLD}  ║  Stress: 600 requests across 6 endpoints             ║${NC}"
echo -e "${GREEN}${BOLD}  ╚══════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  ${DIM}Check 'docker compose logs app' for full OTEL hook traces${NC}"
echo ""
