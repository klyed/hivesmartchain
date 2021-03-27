#!/usr/bin/env bash
# ----------------------------------------------------------
# PURPOSE

# This is the test runner for hsc integration tests.
# It is responsible for starting up a single node hsc test chain and tearing it down afterwards.

# ----------------------------------------------------------
# REQUIREMENTS

# * GNU parallel
# * jq

# ----------------------------------------------------------
# USAGE
# source test_runner.sh

script_dir="$( cd "$( dirname "${BASH_SOURCE}" )" && pwd )"

export hsc_bin="./bin/hsc" #${hsc_bin:-hsc} :-/mnt/c/users/klye/projects/hivesmartchain/bin/hsc
export solc_bin="./bin/solc" #:-/mnt/c/users/klye/projects/hivesmartchain/tests/scripts/deps/solang/solang-linux
export solang_bin="./bin/solang" #:-/mnt/c/users/klye/projects/hivesmartchain/tests/scripts/deps/solc/solc

# If false we will not try to start HiveSmartChain and expect them to be running
export boot=${boot:-true}
export debug=${debug:-true}
export clean=${clean:-true}

export test_exit=0

if [[ "$debug" = true ]]; then
    set -o xtrace
fi


# Note: do not set -e in order to capture exit correctly in mocha
# ----------------------------------------------------------
# Constants

# Ports etc must match those in hsc.toml
export HSC_HOST=127.0.0.1
export HSC_GRPC_PORT=10997


export chain_dir="$script_dir/chain"
export hsc_root="$chain_dir/.testnet"

# Temporary logs
export hsc_log="$chain_dir/hsc.log"
#
# ----------------------------------------------------------

# ---------------------------------------------------------------------------
# Needed functionality

pubkey_of() {
    jq -r ".Accounts | map(select(.Name == \"$1\"))[0].PublicKey.PublicKey" "$chain_dir/genesis.json"
}

address_of() {
    jq -r ".Accounts | map(select(.Name == \"$1\"))[0].Address" "$chain_dir/genesis.json"
}

test_setup(){
  echo "Setting up..."
  cd "$script_dir"

  echo
  echo "Using binaries:"
  echo "  $(type ${solc_bin}) (version: $(./${solc_bin} --version))"
  echo "  $(type ${solang_bin}) (version: $(./${solang_bin} --version))"
  echo "  $(type ${hsc_bin}) (version: $(./${hsc_bin} --version))"
  echo
  # start test chain
  HSC_ADDRESS="$HSC_HOST:$HSC_GRPC_PORT"
  if [[ "$boot" = true ]]; then
    echo "Starting HiveSmartChain using GRPC address: $HSC_ADDRESS..."
    echo
    rm -rf ${hsc_root}
    pushd "$chain_dir"
    ${hsc_bin} start --index 0 --grpc-address $HSC_ADDRESS 2> "$hsc_log"&
    hsc_pid=$!
    popd
  else
    echo "Not booting HiveSmartChain, but expecting HiveSmartChain to be running with tm RPC on port $HSC_GRPC_PORT"
  fi

  export key1_addr=$(address_of "Full_0")
  export key2_addr=$(address_of "Participant_0")
  export key1=Full_0
  export key2=Participant_0
  export key2_pub=$(pubkey_of "Participant_0")

  echo -e "Default Key =>\t\t\t\t$key1_addr"
  echo -e "Backup Key =>\t\t\t\t$key2_addr"
  sleep 4 # boot time

  echo "Setup complete"
  echo ""
}

test_teardown(){
  echo "Cleaning up..."
  if [[ "$boot" = true ]]; then
    echo "Killing hsc with PID $hsc_pid"
    kill ${hsc_pid} 2> /dev/null
    echo "Waiting for hsc to shutdown..."
    wait ${hsc_pid} 2> /dev/null
    rm -rf "$hsc_root"
  fi
  echo ""
  if [[ "$test_exit" -eq 0 ]]
  then
    [[ "$boot" = true ]] && rm -f "$hsc_log"
    echo "Tests complete! Tests are Green. :)"
  else
    echo "Tests complete. Tests are Red. :("
   fi
  exit ${test_exit}
}
