#!/usr/bin/env bash
set -e
# Static solc that will run on linux included Alpine
SOLC_URL="https://github.com/ethereum/solidity/releases/download/v0.5.12/solc-static-linux"
SOLC_BIN=""./mnt/c/users/klye/projects/hivesmartchain/tests/scripts/deps/solc/solc"

#wget -O "$SOLC_BIN" "$SOLC_URL"

sudo chmod +x "$SOLC_BIN"
