#!/usr/bin/env bash
set -e
SOLANG_URL="https://github.com/hyperledger-labs/solang/releases/download/v0.1.7/solang-linux"
SOLANG_BIN="./mnt/c/users/klye/projects/hivesmartchain/tests/scripts/deps/solang/solang"

#wget -O "$SOLANG_BIN" "$SOLANG_URL"

chmod +x "$SOLANG_BIN"
