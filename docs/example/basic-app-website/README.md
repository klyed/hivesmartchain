# Example App (& Website)

Note that this example is the same as the the basic-app example, apart from step 3.

This example contains an example solidity contract [simplestorage](simplestorage.sol) and a [node.js application](app.js) that interacts with the contract using [hsc](../../../js/README.md). It also contains a [makefile](makefile) that will set up a single node chain, deploy the contract using `hsc deploy`. The node app configures itself to use the the single node chain my looking for [account.json](account.json) and [deploy.output.json](deploy.output.json) files that are emitted by `hsc deploy` and the makefile.

The makefile provides some examples of using the `hsc` command line tooling and you are invited to modify it for your purposes (i.e. change from linux target to darwin)

## Dependencies
To run the makefile you will need to have installed:

- GNU Make
- Node.js (the `node` binary)
- npm (the node package manager)
- jq (the JSON tool)
- GO
- Solc (solidity compiler)

HiveSmartChain will be downloaded for you when using the makefile, but you may override `HSC_BIN` and `HSC_ARCH` in the makefile to change this behaviour. By default HiveSmartChain is downloaded for `Linux_x86_64.

## Running the example

All commands should be run from the same directory as this readme file.

### Step one
Start the chain

```shell
make start_chain
```

This will install hsc, create a new chain as necessary.

If successful you will see continuous output in your terminal, you can shutdown HiveSmartChain by sending the interrupt signal with Ctrl-C, and restart it again with whatever state has accumulated with `make start_chain`. If you would like to destroy the existing chain and start completely fresh (including deleting keys) run `make rechain`. If you would like to keep existing keys and chain config run `make reset_chain`.

You can redeploy the contract (to a new address) with `make redeploy`. The node app will then use this new contract by reading the address deploy.output.json. Be sure to do this if you wish to modify simplestorage.sol.

### Step two
Leave hsc running and in a separate terminal start the app which runs a simple HTTP server with:

```shell
make start_app
```

This will deploy the contract if necessary, install any node dependencies, and then start an expressjs server, which will run until interrupted.

### Step three
Open a web browser and type:

```shell
  http://localhost:3000/
```

You will see two buttons:
* Set Value - that allows you to change the value stored in the associated smart contract
* Get Value - that allows you to retrieve the value stored in the associated smart contract, which will be displayed underneath the button

