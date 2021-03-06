package bcm

import (
	"github.com/klyed/hivesmartchain/txs"
	"github.com/tendermint/tendermint/types"
)

type Block struct {
	txDecoder txs.Decoder
	*types.Block
}

func NewBlock(txDecoder txs.Decoder, block *types.Block) *Block {
	return &Block{
		txDecoder: txDecoder,
		Block:     block,
	}
}

func (b *Block) Transactions(iter func(*txs.Envelope) error) error {
	for i := 0; i < len(b.Txs); i++ {
		tx, err := b.txDecoder.DecodeTx(b.Txs[i])
		if err != nil {
			return err
		}
		err = iter(tx)
		if err != nil {
			return err
		}
	}
	return nil
}
