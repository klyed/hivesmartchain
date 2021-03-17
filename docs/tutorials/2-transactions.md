# Transactions

Burrow supports a number of [transactions](reference/transactions.md) which denote a unit of computation.
The easiest way to experiment is with our `hsc tx` command, but please checkout the [deployment guide](deploy.md)
for more advanced usage.

## Getting Started

Let's start a chain with one validator to process blocks and two participant accounts:

```shell
hsc spec -v1 -p2 | hsc configure -s- > hsc.toml
hsc start -v0 &
```

Make a note of the two participant addresses generated in the `hsc.toml`.

## Send Token

Let's formulate a transaction to send funds from one account to another.
Given our two addresses created above, set `$SENDER` and `$RECIPIENT` respectively.
We'll also need to designate an amount of native token available from our sender.

```shell
hsc tx formulate send -s $SENDER -t $RECIPIENT -a $AMOUNT > tx.json
```

To send this transaction to your local node and subsequently the chain (if running more than one validator),
pipe the output above through the following command:

```shell
hsc tx commit --file tx.json
```