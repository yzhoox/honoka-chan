#!/usr/bin/env bash
set -euo pipefail

source "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)/common.sh"

require_cmd zip

ZIP_PATH="$ANDROID_DIR/app/src/main/assets/honoka_runtime.zip"

cd "$PROJECT_ROOT"

ensure_file "$PROJECT_ROOT/assets/main.db"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

mkdir -p "$TMP_DIR/assets" "$TMP_DIR/static"

log "collecting runtime files"

cat > "$TMP_DIR/config.json" <<'EOF'
{
	"app_name": "honoka-chan",
	"settings": {
		"listen_port": "8080",
		"cdn_server": "http://127.0.0.1:8080/static",
		"reload_token": "",
		"unlock_all_special_rotation": false
	}
}
EOF

cp -r "$PROJECT_ROOT/assets/main.db" "$TMP_DIR/assets/main.db"

for path in \
  "assets/serverdata" \
  "assets/certs" \
  "static/templates" \
  "static/css" \
  "static/js" \
  "static/images" \
  "static/font"
do
  if [[ -e "$PROJECT_ROOT/$path" ]]; then
    mkdir -p "$TMP_DIR/$(dirname "$path")"
    cp -r "$PROJECT_ROOT/$path" "$TMP_DIR/$path"
  fi
done

rm -f "$ZIP_PATH"

log "packing $ZIP_PATH"
(
  cd "$TMP_DIR"
  zip -qr "$ZIP_PATH" assets static config.json
)

log "runtime zip ready: $ZIP_PATH"
