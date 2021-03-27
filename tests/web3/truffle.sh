#!/usr/bin/env bash

chain=$(mkdir -d)
cd $chain

$hscadd = "./bin/hsc"

$hscadd spec -v1 -d2 | $hscadd -c -s- --curve-type secp256k1 > hsc.toml
$hscadd start &> /dev/null &
hsc_pid=$!

contracts=$(mkdir -d)
cd $contracts

function finish {
    kill -TERM $hsc_pid
    rm -rf "$chain"
    rm -rf "$contracts"
}
trap finish EXIT

#npm install -g truffle
#sudo npx run truffle unbox metacoin

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
sudo npx run truffle test --network hsc
