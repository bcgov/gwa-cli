#!/usr/bin/env bash
set -euo pipefail

# Load .env if present
if [ -f ".env" ]; then
  set -a
  source .env
  set +a
fi

# ---- adjust these ----
GWA_BIN="${GWA_BIN:-./gwa}"
GWA_CONFIG="${GWA_CONFIG:-$HOME/.gwa-config.yaml}"
DEV_HOST="${GWA_API_HOST:-api-gov-bc-ca.dev.api.gov.bc.ca}"
SCHEME="${SCHEME:-https}"
CLIENT_ID="${GWA_CLIENT_ID:-gwa-cli}"
# ----------------------

timestamp() {
  date +"%Y-%m-%d %H:%M:%S"
}

log() {
  printf '[%s] %s\n' "$(timestamp)" "$*"
}

fail() {
  log "FAIL: $*"
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing required command: $1" >&2
    exit 1
  }
}

backup_config() {
  if [ -f "$GWA_CONFIG" ]; then
    cp "$GWA_CONFIG" "$GWA_CONFIG.bak.pkce-test"
    log "Backed up config to $GWA_CONFIG.bak.pkce-test"
  fi
}

restore_config() {
  if [ -f "$GWA_CONFIG.bak.pkce-test" ]; then
    mv -f "$GWA_CONFIG.bak.pkce-test" "$GWA_CONFIG"
    log "Restored original config"
  fi
}

write_config_base() {
  mkdir -p "$(dirname "$GWA_CONFIG")"
  cat > "$GWA_CONFIG" <<EOF
host: $DEV_HOST
scheme: $SCHEME
EOF
}

set_pkce_method() {
  local method="$1"

  write_config_base

  if [ -n "$method" ]; then
    cat >> "$GWA_CONFIG" <<EOF
pkce_method: $method
EOF
  fi
}

run_case() {
  local label="$1"
  local method="$2"
  local expect="$3"
  local rc

  log "============================================================"
  log "CASE: $label"
  log "pkce_method: ${method:-<unset>}"
  log "expected: $expect"

  set_pkce_method "$method"

  log "Config written to $GWA_CONFIG:"
  sed 's/^/  /' "$GWA_CONFIG"

  set +e
  "$GWA_BIN" login
  rc=$?
  set -e

  log "exit code: $rc"

  case "$expect" in
    fail)
      if [ "$rc" -eq 0 ]; then
        fail "Case '$label' was expected to fail, but it succeeded. Non-S256 PKCE behavior is more permissive than expected."
      fi
      log "Observed failure as expected"
      ;;
    interactive-success)
      if [ "$rc" -ne 0 ]; then
        fail "Case '$label' was expected to succeed, but it failed with exit code $rc."
      fi
      log "Observed success"
      ;;
    *)
      fail "Unknown expectation: $expect"
      ;;
  esac
}

cleanup() {
  restore_config || true
}
trap cleanup EXIT

require_cmd sed

if [ ! -x "$GWA_BIN" ]; then
  fail "GWA binary not found or not executable: $GWA_BIN. Build it first."
fi

backup_config

log "Starting PKCE device-login test matrix"
log "This script enforces the target behavior:"
log "  unset -> success (CLI defaults to S256)"
log "  plain -> fail"
log "  S256 -> success"
log "  None -> fail"
log "  S512 -> fail"

run_case "unset pkce_method (CLI default path)" "" "interactive-success"
run_case "explicit plain" "plain" "fail"
run_case "explicit S256" "S256" "interactive-success"
run_case "explicit None" "None" "fail"
run_case "explicit S512" "S512" "fail"

log "PKCE test matrix complete"