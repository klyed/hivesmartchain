# Basics

You can spin up a single node chain with:

```shell
hsc spec -v1 | hsc configure -s- | hsc start -c-
```

## Configuration

The quick-and-dirty one-liner looks like:

```shell
# Read spec on stdin
hsc spec -r1 -p10 -f1 | hsc configure -s- > hsc.toml
```

Which translates into:

```shell
hsc spec --participant-accounts=1 --full-accounts=1 > genesis-spec.json
hsc configure --genesis-spec=genesis-spec.json > hsc.toml
```

> You might want to run this in a clean directory to avoid overwriting any previous spec or config.

## Running

Once the `hsc.toml` has been created, we run:

```
# To select our validator address by index in the GenesisDoc
hsc start --validator=0
# Or to select based on address directly (substituting the example address below with your validator's):
hsc start --address=BE584820DC904A55449D7EB0C97607B40224B96E
```

If you would like to reset your node, you can just delete its working directory with `rm -rf .hsc`.
In the context of a multi-node chain it will resync with peers, otherwise it will restart from height 0.

## Keys

HiveSmartChain consumes its keys through our key signing interface that can be run as a standalone service with:

```shell
hsc keys server
```

This command starts a key signing daemon capable of generating new ed25519 and secp256k1 keys, naming those keys, signing arbitrary messages, and verifying signed messages.
It also initializes a key store directory in `.keys` (by default) where private key matter is stored.

It should be noted that the GRPC service exposed by the keys server will sign _any_ inbound requests using the keys it maintains so the machine running the keys service should only allow connections from sources that are trusted to use those keys.
