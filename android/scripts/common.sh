#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
ANDROID_DIR="$(cd -- "$SCRIPT_DIR/.." && pwd)"
PROJECT_ROOT="$(cd -- "$ANDROID_DIR/.." && pwd)"

readonly SCRIPT_DIR
readonly ANDROID_DIR
readonly PROJECT_ROOT

log() {
  printf '[android] %s\n' "$*"
}

die() {
  printf '[android] error: %s\n' "$*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing command: $1"
}

ensure_file() {
  [[ -f "$1" ]] || die "missing file: $1"
}
