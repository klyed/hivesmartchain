package exec

import (
	"testing"

	"github.com/klyed/hivesmartchain/event/query"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestBlockExecution_Marshal(t *testing.T) {
	be := &BlockExecution{
		Header: &tmproto.Header{
			Height:          3,
			AppHash:         []byte{2},
			ProposerAddress: []byte{1, 2, 33},
		},
	}
	bs, err := be.Marshal()
	require.NoError(t, err)
	beOut := new(BlockExecution)
	require.NoError(t, beOut.Unmarshal(bs))
}

func TestBlockExecution_StreamEvents(t *testing.T) {
	be := &BlockExecution{
		Header: &tmproto.Header{
			Height:          2,
			AppHash:         []byte{2},
			ProposerAddress: []byte{1, 2, 33},
		},
	}

	qry, err := query.NewBuilder().AndContains("Height", "2").Query()
	require.NoError(t, err)
	match := qry.Matches(be)
	require.True(t, match)
}
