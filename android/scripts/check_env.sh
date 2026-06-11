#!/usr/bin/env bash
set -euo pipefail

source "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)/common.sh"

ok() {
  printf '[android] ok: %s\n' "$*"
}

warn() {
  printf '[android] warn: %s\n' "$*" >&2
}

status=0

check_cmd() {
  local cmd="$1"
  if command -v "$cmd" >/dev/null 2>&1; then
    ok "command found: $cmd -> $(command -v "$cmd")"
  else
    warn "command missing: $cmd"
    status=1
  fi
}

check_file() {
  local file="$1"
  if [[ -f "$file" ]]; then
    ok "file found: $file"
  else
    warn "file missing: $file"
    status=1
  fi
}

check_dir() {
  local dir="$1"
  if [[ -d "$dir" ]]; then
    ok "directory found: $dir"
  else
    warn "directory missing: $dir"
    status=1
  fi
}

log "checking common commands"
check_cmd go
check_cmd zip

SDK_ROOT="${ANDROID_SDK_ROOT:-${ANDROID_HOME:-}}"
NDK_ROOT="${ANDROID_NDK_HOME:-${ANDROID_NDK_ROOT:-}}"
HOST_TAG="${ANDROID_NDK_HOST_TAG:-linux-x86_64}"
MIN_SDK="${MIN_SDK:-26}"

if [[ -n "$SDK_ROOT" ]]; then
  check_dir "$SDK_ROOT"
else
  warn "ANDROID_SDK_ROOT is not set"
  status=1
fi

if [[ -n "$NDK_ROOT" ]]; then
  check_dir "$NDK_ROOT"
  TOOLCHAIN="$NDK_ROOT/toolchains/llvm/prebuilt/$HOST_TAG/bin"
  check_dir "$TOOLCHAIN"
  check_file "$TOOLCHAIN/aarch64-linux-android${MIN_SDK}-clang"
  check_file "$TOOLCHAIN/armv7a-linux-androideabi${MIN_SDK}-clang"
  check_file "$TOOLCHAIN/x86_64-linux-android${MIN_SDK}-clang"
else
  warn "ANDROID_NDK_HOME is not set"
  status=1
fi

log "checking project files"
check_file "$PROJECT_ROOT/assets/main.db"
check_dir "$PROJECT_ROOT/assets/serverdata"
check_dir "$PROJECT_ROOT/assets/certs"
check_dir "$PROJECT_ROOT/static/templates"

if [[ $status -ne 0 ]]; then
  warn "environment check failed"
  exit $status
fi

ok "environment check passed"
