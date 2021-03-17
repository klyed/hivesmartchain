#!/usr/bin/env bash

chain=$(mktemp -d)
cd $chain
$hsc_bin spec -v1 -d2 | $hsc_bin configure -s- --curve-type secp256k1 > hsc.toml
$hsc_bin start &> /dev/null &
hsc_pid=$!

contracts=$(mktemp -d)
cd $contracts

function finish {
    kill -TERM $hsc_pid
    rm -rf "$chain"
    rm -rf "$contracts"
}
trap finish EXIT

npm install -g truffle
truffle unbox metacoin

cat << EOF > truffle-config.js
module.exports = {
  networks: {
   hsc: {
     host: "127.0.0.1",
     port: 26660,
     network_id: "*",
   },
  }
};
EOF
truffle test --network hsc
