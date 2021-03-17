# [Hive Smart Chain (HSC)](https://peakd.com/proposals/164)
## Bringing Truly Decentralized Smart Contracts to HIVE

<!--
[![CI](https://github.com/klyed/hivesmartchain/workflows/main/badge.svg)](https://launch-editor.github.com/actions?workflowID=main&event=push&nwo=hyperledger%2Fhsc)
[![version](https://img.shields.io/github/tag/hyperledger/hsc.svg)](https://github.com/klyed/hivesmartchain/releases/latest)
[![GoDoc](https://godoc.org/github.com/hsc?status.png)](https://godoc.org/github.com/hyperledger/hsc)
[![license](https://img.shields.io/github/license/hyperledger/hsc.svg)](../LICENSE.md)
[![LoC](https://tokei.rs/b1/github/hyperledger/hsc?category=lines)](https://github.com/hyperledger/hsc)
[![codecov](https://codecov.io/gh/hyperledger/hsc/branch/main/graph/badge.svg)](https://codecov.io/gh/hyperledger/hsc)
-->

Hive Smart Chain is a community focused Ethereum smart-contract capable side chain node allowing HIVE users to run smart contracts while using HIVE as their gas fees. It executes Ethereum EVM and WASM smart contract code (usually written in [Solidity](https://solidity.readthedocs.io)) on a permissioned virtual machine. Hive Smart Chain provides transaction finality and high transaction throughput on a proof-of-stake [Tendermint](https://tendermint.com) consensus engine.

![hsc logo](https://images.hive.blog/0x0/https://files.peakd.com/file/peakd-hive/klye/23uQERcpCeRtQUooESeioNqouBjmc9mDBsT68U4AwRDdraYPzwa4sDsSLpVNt4ohW3uun.png)

## What is Hive Smart Chain

Hive Smart Chain is a fully fledged blockchain node and smart contract execution engine that runs in tandem to the HIVE DPoS blockchain offering it access to smart contract capabilities for it's native token as well as acting as a gateway to interface with other blockchains in a "DeFi" manner -- in otherwords it's a distributed database and computing platform that executes code. Hive Smart Chain runs Ethereum Virtual Machine (EVM) and Web Assembly (WASM) smart contracts. Hive Smart Chain are Secured and decentralized using the [Tendermint](https://github.com/tendermint/tendermint) consensus algorithm.

Highlights include:

- **Tamper-resistant merkle state** - a node can detect if its state is corrupted or if a validator is dishonestly executing the protocol.
- **Proof-of-stake support** - run a private or public permissioned network.
- **On-chain governance primitives** - stakeholders may vote for autonomous smart contract upgrades.
- **Ethereum account world-view** - state and code is organised into cryptographically-addressable accounts.
- **Low-level permissioning** - code execution permissions can be set on a per-account basis.
- **Event streaming** - application state is organised in an event stream and can drive external systems.
- **[SQL mapping layer](reference/vent.md)** - map smart contract event emissions to SQL tables using a projection specification.
- **GRPC interfaces** - all RPC calls can be accessed from any language with GRPC support. Protobuf is used extensively for core objects.
- **Javascript client library** - client library uses code generation to provide access to contracts via statically Typescript objects.
- **Keys service** - provides optional delegating signing at the server or via a local proxy
- **Web3 RPC** - provides compatibility for mainnet Ethereum tooling such as Truffle and Metamask

## JavaScript Client

There is a [JavaScript API](https://github.com/klyed/hivesmartchain/tree/main/js)

## Useful commands

###start node / single user / test mode
./hsc spec -f1 | ./hsc configure -s- | ./hsc start -v0 -c-

## Project Roadmap

Project information generally updated on a weekly basis can be found on the [Hive Smart Chain](https://github.com/klyed/hivesmartchain).


## Releases

- **Hive Smart Chain binaries**: https://github.com/klyed/hivesmartchain/releases
- **hsc.js   (COMING SOON)**: https://www.npmjs.com/package/@hyperledger/hsc
- **Docker   (COMING SOON)**: https://hub.docker.com/repository/docker/hyperledger/hsc

## Contribute

We welcome any and all contributions. Read the [contributing file](../.github/CONTRIBUTING.md) for more information on making your first Pull Request to HiveSmartChain!

You can find us on:
- [PeakD](https://peakd.com/@klye)
- [Discord](https://lists.hyperledger.org/mailman/listinfo)
- [here on Github](https://github.com/klyed/hivesmartchain/issues)

## License

[Apache 2.0](../LICENSE.md)
