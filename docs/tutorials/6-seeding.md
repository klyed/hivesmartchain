# Seed Nodes

## What is a seed node?

According to [Tendermint documentation](https://tendermint.com/docs/tendermint-core/using-tendermint.html#seed):
>A seed node is a node who relays the addresses of other peers which they know
of. These nodes constantly crawl the network to try to get more peers. The
addresses which the seed node relays get saved into a local address book. Once
these are in the address book, you will connect to those addresses directly.
Basically the seed nodes job is just to relay everyones addresses. You won't
connect to seed nodes once you have received enough addresses, so typically you
only need them on the first start. The seed node will immediately disconnect
from you after sending you some addresses.

### Seed Mode

If a node is in seed mode it will accept inbound connections, share its address book, then hang up.
Seeds modes will do a bit of gossip but not that usefully.
Any type of node can be referenced as a 'Seeds' in the config, whether or not another node considers this node as a seed is independent of whether this node is in 'seed mode'.
These are different concepts:
> You are free to use a non-seed-mode node as a seed.

You do not want to have too many seeds in your network (because they just keep hanging up on other nodes once they've spread their wild oats), but they are useful for accelerating peer exchange (of addresses).

### Persistent Peers
Persistent peers are peers that you want to connect of regardless of the heuristics and churn dynamics built into the p2p switch.
Ordinarily you would not stay connected to a particular peer forever, and you would not indefinitely redial, but you will for a persistent peer.

## Configure

In this quick start, we will create validator nodes which do not know about each other.
A seed node will crawl the network and relay addresses.

### Seed Node

```shell
hsc spec -f1 | hsc configure --keys-dir=.keys_seed -s- > /dev/null
```

```toml
HscDir = ".hsc_seed_0"

[Tendermint]
  SeedMode = true
  ListenHost = "0.0.0.0"
  ListenPort = "10000"
  Moniker = "seed_node_0"

[Execution]

[Keys]
  GRPCServiceEnabled = false
  AllowBadFilePermissions = true
  RemoteAddress = ""
  KeysDirectory = ".keys_seed"

[RPC]
  [RPC.Info]
    Enabled = true
    ListenHost = "127.0.0.1"
    ListenPort = "10001"
  [RPC.Profiler]
    Enabled = false
  [RPC.GRPC]
    Enabled = false
  [RPC.Metrics]
    Enabled = false
```

### Validators

```shell
hsc spec --full-accounts=3 | hsc configure -s- > .hsc_init.toml
```

From the generated `.hsc_init.toml` file, create new files for each node, and change the content.

#### Validator 1

```toml
HscDir = ".hsc_node0"

[Tendermint]
  Seeds = "PUT_HERE_SEED_NODE_ID@LISTEN_EXTERNAL_ADDRESS"
  SeedMode = false
  PersistentPeers = ""
  ListenHost = "0.0.0.0"
  ListenPort = "20000"
  Moniker = "val_node_0"

[Execution]

[Keys]
  GRPCServiceEnabled = false
  AllowBadFilePermissions = true
  RemoteAddress = ""
  KeysDirectory = ".keys"

[RPC]
  [RPC.Info]
    Enabled = true
    ListenHost = "127.0.0.1"
    ListenPort = "20001"
  [RPC.Profiler]
    Enabled = false
  [RPC.GRPC]
    Enabled = true
    ListenHost = "127.0.0.1"
    ListenPort = "20002"
  [RPC.Metrics]
    Enabled = false
```

#### Validator 2

```toml
HscDir = ".hsc_node1"

[Tendermint]
  Seeds = "PUT_HERE_SEED_NODE_ID@LISTEN_EXTERNAL_ADDRESS"
  SeedMode = false
  PersistentPeers = ""
  ListenHost = "0.0.0.0"
  ListenPort = "30000"
  Moniker = "val_node_1"

[Execution]

[Keys]
  GRPCServiceEnabled = false
  AllowBadFilePermissions = true
  RemoteAddress = ""
  KeysDirectory = ".keys"

[RPC]
  [RPC.Info]
    Enabled = true
    ListenHost = "127.0.0.1"
    ListenPort = "30001"
  [RPC.Profiler]
    Enabled = false
  [RPC.GRPC]
    Enabled = true
    ListenHost = "127.0.0.1"
    ListenPort = "30002"
  [RPC.Metrics]
    Enabled = false
```

#### Validator 3

```toml
HscDir = ".hsc_node2"

[Tendermint]
  Seeds = "PUT_HERE_SEED_NODE_ID@LISTEN_EXTERNAL_ADDRESS"
  SeedMode = false
  PersistentPeers = ""
  ListenHost = "0.0.0.0"
  ListenPort = "40000"
  Moniker = "val_node_2"

[Execution]

[Keys]
  GRPCServiceEnabled = false
  AllowBadFilePermissions = true
  RemoteAddress = ""
  KeysDirectory = ".keys"

[RPC]
  [RPC.Info]
    Enabled = true
    ListenHost = "127.0.0.1"
    ListenPort = "40001"
  [RPC.Profiler]
    Enabled = false
  [RPC.GRPC]
    Enabled = true
    ListenHost = "127.0.0.1"
    ListenPort = "40002"
  [RPC.Metrics]
    Enabled = false
```

## Start Network

### Seed Node

```shell
hsc start --address=`basename .keys_seed/data/* .json` --config=.hsc_seed.toml  > .hsc_seed.log 2>&1 &
```

#### Validators

Tendermint requires strict and routable address (not loopback, local etc), you can find the listen address with this command:

```shell
SEED_URL=`curl -s 127.0.0.1:10001/network | jq -r '.result.ThisNode | [.ID, .ListenAddress] | join("@") | ascii_downcase'`
echo $SEED_URL
```

Configure the validator nodes to connect to the seed node:

```shell
sed -i s%PUT_HERE_SEED_NODE_ID@LISTEN_EXTERNAL_ADDRESS%${SEED_URL}% .hsc_val0.toml
sed -i s%PUT_HERE_SEED_NODE_ID@LISTEN_EXTERNAL_ADDRESS%${SEED_URL}% .hsc_val1.toml
sed -i s%PUT_HERE_SEED_NODE_ID@LISTEN_EXTERNAL_ADDRESS%${SEED_URL}% .hsc_val2.toml
```

Start the network:

```shell
hsc start -v=0 --config=.hsc_val0.toml  > .hsc_val0.log 2>&1 &
hsc start -v=1 --config=.hsc_val1.toml  > .hsc_val1.log 2>&1 &
hsc start -v=2 --config=.hsc_val2.toml  > .hsc_val2.log 2>&1 &
```

The nodes should connect to our seed node and request addresses, then they will connect to each other and start submitting and voting on blocks.


To check the network status, and that the validator nodes are connected to each other run:

```shell
curl -s 127.0.0.1:40001/network | jq -r '.result.peers[].node_info.moniker'
val_node_0
val_node_1
```

You can monitor consensus and current blockchain height from the node info websocket:

```shell
curl -s 127.0.0.1:20001/consensus | jq -r '.result.round_state.height'
```

Disable seed mode on the seed node and see how it affects the peers network:

```toml
[Tendermint]
  SeedMode = false
```

Clear nodes folder (it will restart the chain from the genesis block):

```shell
killall hsc
rm -rf .hsc_node0 .hsc_node1 .hsc_node2 .hsc_seed_0
```

Restart all nodes, then check network status (Validator 3 is now connected to all peers, included seed node):

```shell
curl -s 127.0.0.1:40001/network | jq -r '.result.peers[].node_info.moniker'
seed_node_0
val_node_0
val_node_1
```
