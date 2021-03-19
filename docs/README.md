# [Hive Smart Chain (HSC)](https://peakd.com/proposals/164)
## Bringing Truly Decentralized Smart Contracts to HIVE Integrating 
## the Ability to Interact with Industry Standard Contract Execution

<!--
[![CI](https://github.com/klyed/hivesmartchain/workflows/main/badge.svg)](https://launch-editor.github.com/actions?workflowID=main&event=push&nwo=hyperledger%2Fhsc)
[![version](https://img.shields.io/github/tag/hyperledger/hsc.svg)](https://github.com/klyed/hivesmartchain/releases/latest)
[![GoDoc](https://godoc.org/github.com/hsc?status.png)](https://godoc.org/github.com/hyperledger/hsc)
[![license](https://img.shields.io/github/license/hyperledger/hsc.svg)](../LICENSE.md)
[![LoC](https://tokei.rs/b1/github/hyperledger/hsc?category=lines)](https://github.com/hyperledger/hsc)
[![codecov](https://codecov.io/gh/hyperledger/hsc/branch/main/graph/badge.svg)](https://codecov.io/gh/hyperledger/hsc)
-->

Hive Smart Chain is a community focused Ethereum smart-contract capable side chain node allowing HIVE users to run smart contracts while using HIVE as their gas fees. It executes Ethereum EVM and WASM smart contract code (usually written in [Solidity](https://solidity.readthedocs.io)) on a permissioned virtual machine. Hive Smart Chain provides transaction finality and high transaction throughput on a proof-of-stake [Tendermint](https://tendermint.com) consensus engine.

The Hive Smart Chain will not only serve as a great addition of industry standard smart conrtract typing introducing proper ERC-20 and ERC-271 NFT but also as a means of bringing the HIVE blockchain into the world of decentralized finance with the ability to excersize great blockchain interoperability due to it's modular nature.

![hsc logo](https://images.hive.blog/0x0/https://files.peakd.com/file/peakd-hive/klye/23uQERcpCeRtQUooESeioNqouBjmc9mDBsT68U4AwRDdraYPzwa4sDsSLpVNt4ohW3uun.png)

## What is Hive Smart Chain

Hive Smart Chain is a fully fledged blockchain node and smart contract execution engine that runs in tandem to the HIVE DPoS blockchain offering it access to smart contract capabilities for it's native token as well as acting as a gateway to interface with other blockchains in a "DeFi" manner -- in otherwords it's a distributed database and computing platform that executes code. Hive Smart Chain runs Ethereum Virtual Machine (EVM) and Web Assembly (WASM) smart contracts. Hive Smart Chain are Secured and decentralized using the [Tendermint](https://github.com/klyed/tendermint) consensus algorithm.

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



# Commands and Capabilities

## Transactions

HiveSmartChain implements a number of transaction types. Transactions will be ordered by our consensus mechanism (Tendermint) and applied to our application state machine - 
replicated across all HiveSmartChain nodes. Each transaction is applied atomically and runs deterministically. The transactions contain the arguments for an 
[execution context](https://github.com/klyed/hivesmartchain/tree/main/execution/contexts).

Our transactions are defined in Protobuf [here](https://github.com/klyed/hivesmartchain/blob/main/protobuf/payload.proto).

Transactions can be built using our GRPC client libraries programmatically, via [hsc.js](js-api.md), or with `hsc deploy` - see our [deployment guide](deploy.md).

## TxInput

| Parameter | Type | Description |
| ----------|------|-------------|
| Address | Address | The address of an account issuing this transaction - the transaction envelope must also be signed by the private key associated with this address |
| Amount | uint64 | The amount of native token to transfer from the input to the output of the transaction |
| Sequence | uint64 | A counter that must match the current value of the input account's Sequence plus one - i.e. the Sequence must equal n if this is the nth transaction issued by this account |


## CallTx

Our core transaction type that calls EVM code, possibly transferring value. It takes the following parameters:

| Parameter | Type | Description |
| ----------|------|-------------|
| Input | TxInput | The external 'caller' account - will be the initial SENDER and CALLER |
| Address | *Address | The address 'callee' contract - the contract whose code will be executed. If this value is nil then the CallTx is interpreted as contract creation and will deploy the bytecode contained in Data or WASM |
| GasLimit | uint64 | The maximum number of computational steps that we will allow to run before aborted the transaction execution. Measured according to our hardcoded simplified gas schedule (one gas unit per operation). Ensure transaction termination. If 0 a default cap will be used. |
| Fee | uint64 | An optional fee to be subtracted from the input amount - currently this fee is simply burnt! In the future fees will be collected and disbursed amongst validators as part of our token economics system |
| Data | []byte |  If the CallTx is a deployment (i.e. Address is nil) then this data will be executed as EVM bytecode will and the return value will be used to instatiate a new contract. If the CallTx is a plain call then the data will form the input tape for the EVM call |

## SendTx

Allows [native token](reference/participants.md) to be sent from multiple inputs to multiple outputs. The basic value transfer function that calls no EVM Code.

## NameTx

Provides access to a global name registry service that associates a particular string key with a data payload and an owner. The control of the name is guaranteed for 
the period of the lease which is a determined by a fee.

> A future revision will change the way in which leases are calculated. Currently we use a somewhat historically-rooted fixed fee, see the [`NameCostPerBlock` function](https://github.com/klyed/hivesmartchain/blob/main/execution/names/names.go#L83).

## BondTx

This allows validators nominate themselves to the validator set by placing a bond subtracted from their balance.

For more information see the [bonding documentation](reference/bonding.md).

## UnbondTx

This allows validators remove themselves to the validator set returning their bond to their balance.

## BatchTx

Runs a set of transactions atomically in a single meta-transaction within a single block

## GovTx

An all-powerful transaction for modifying existing accounts.

## ProposalTx

A transaction type containing a batch of transactions on which a ballot is held to determine whether to execute, see the [proposals tutorial](tutorials/8-proposals.md).

## PermsTx

A transaction to modify the permissions of accounts.

## IdentifyTx

When running a closed or permissioned network, it is desirable to restrict the participants.
For example, a consortium may wish to run a shared instance over a wide-area network without
sharing the state to unknown parties. 

As Tendermint handles P2P connectivity for HiveSmartChain, it extends a concept known as the 'peer filter'.
This means that on every connection request to a particular node, our app will receive a request to 
check a whitelist (if enabled, otherwise allowed by default) - if the source IP address or node key is 
unknown then the connection will be rejected. The easiest way to manage this whitelist is to hard code
the respective participants in the config:

```toml
[Tendermint]
  AuthorizedPeers = "DDEF3E93BBF241C737A81E6BA085D0C77C7B51C9@127.0.0.1:26656,
                        A858F15CD7048F7D6C1B310E016A0B8BA1D44861@127.0.0.1:26657"
```

This can become difficult to manage over time, and any change would require a restart of the node. A more
scalable solution is `IdentifyTx`, which allows select participants to be associated with a particular 
'node identity' - network address, node key and moniker. Once enabled in the config, a node will only allow
connection requests from entries in its registry.

```toml
[Tendermint]
  IdentifyPeers = true
```

For more details, see the [ADR](ADRs/adr-2_identify-tx.md).


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
