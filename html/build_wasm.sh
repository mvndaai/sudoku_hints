#!/bin/zsh
# Build Go WASM with RapidAPIKey from config.jsonc (JSON with comments)

CONFIG_FILE="config.jsonc"
# Remove // and /* */ comments, then extract rapidApiKey
KEY=$(grep -vE '^\s*//' "$CONFIG_FILE" | sed 's:/\*.*\*/::g' | jq -r '.rapidApiKey')

go build -tags=js,wasm -ldflags="-X 'main.RapidAPIKey=$KEY'" -o main.wasm
