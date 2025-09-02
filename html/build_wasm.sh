#!/bin/zsh
# Build Go WASM with RapidAPIKey from config.jsonc (JSON with comments)

# Get the directory of this script
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]:-${(%):-%N}}")" && pwd)"

CONFIG_FILE="$SCRIPT_DIR/config.jsonc"
# Remove // and /* */ comments, then extract rapidApiKey
KEY=$(grep -vE '^\s*//' "$CONFIG_FILE" | sed 's:/\*.*\*/::g' | jq -r '.rapidApiKey')
GOOS=js GOARCH=wasm go build -tags=js,wasm -ldflags="-X 'main.RapidAPIKey=$KEY'" -o "$SCRIPT_DIR/main.wasm" "$SCRIPT_DIR/main.go"
