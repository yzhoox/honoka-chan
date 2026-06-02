#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_BIN="$ROOT_DIR/honoka-chan"
PID_FILE="$ROOT_DIR/temp/run/honoka-chan.pid"
LOG_FILE="$ROOT_DIR/temp/run/honoka-chan.log"

mkdir -p "$(dirname "$PID_FILE")"

build_app() {
  if command -v go >/dev/null 2>&1; then
    (cd "$ROOT_DIR" && go build -o honoka-chan main.go)
  else
    echo "Go is not installed, cannot build honoka-chan automatically"
    exit 1
  fi
}

is_running() {
  [[ -f "$PID_FILE" ]] || return 1
  local pid
  pid="$(cat "$PID_FILE" 2>/dev/null || true)"
  [[ -n "${pid:-}" ]] || return 1
  kill -0 "$pid" 2>/dev/null || return 1

  local args
  args="$(ps -p "$pid" -o args= 2>/dev/null || true)"
  [[ "$args" == *"$APP_BIN"* || "$args" == *"honoka-chan"* ]]
}

start_app() {
  if is_running; then
    echo "honoka-chan is already running, pid=$(cat "$PID_FILE")"
    return 0
  fi

  if [[ ! -f "$APP_BIN" ]]; then
    build_app
  fi

  if [[ ! -x "$APP_BIN" ]]; then
    echo "honoka-chan binary is not executable: $APP_BIN"
    exit 1
  fi

  nohup "$APP_BIN" >>"$LOG_FILE" 2>&1 &
  local pid=$!
  echo "$pid" >"$PID_FILE"
  sleep 1

  if kill -0 "$pid" 2>/dev/null; then
    echo "honoka-chan started successfully, pid=$pid"
  else
    rm -f "$PID_FILE"
    echo "honoka-chan failed to start, check $LOG_FILE"
    exit 1
  fi
}

stop_app() {
  if ! is_running; then
    rm -f "$PID_FILE"
    echo "honoka-chan is not running"
    return 0
  fi

  local pid
  pid="$(cat "$PID_FILE")"
  kill "$pid" 2>/dev/null || true

  for _ in $(seq 1 50); do
    if ! kill -0 "$pid" 2>/dev/null; then
      rm -f "$PID_FILE"
      echo "honoka-chan stopped"
      return 0
    fi
    sleep 0.2
  done

  kill -9 "$pid" 2>/dev/null || true
  rm -f "$PID_FILE"
  echo "honoka-chan force stopped"
}

status_app() {
  if is_running; then
    echo "honoka-chan is running, pid=$(cat "$PID_FILE")"
  else
    echo "honoka-chan is not running"
  fi
}

logs_app() {
  touch "$LOG_FILE"
  tail -f "$LOG_FILE"
}

usage() {
  cat <<'EOF'
Usage:
  ./run.sh start    Start service
  ./run.sh stop     Stop service
  ./run.sh restart  Restart service
  ./run.sh status   Check status
  ./run.sh logs     Follow logs
EOF
}

case "${1:-}" in
  start)
    start_app
    ;;
  stop)
    stop_app
    ;;
  restart)
    stop_app
    start_app
    ;;
  status)
    status_app
    ;;
  logs)
    logs_app
    ;;
  ""|-h|--help|help)
    usage
    ;;
  *)
    echo "Unknown command: $1"
    usage
    exit 1
    ;;
esac
