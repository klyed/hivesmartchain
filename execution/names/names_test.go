package names

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeProtobuf(t *testing.T) {
	entry := &Entry{
		Name:    "Foo",
		Data:    "oh noes",
		Expires: 24423432,
		Owner:   crypto.Address{1, 2, 0, 9, 8, 8, 1, 2},
	}
	encoded, err := proto.Marshal(entry)
	require.NoError(t, err)
	entryOut := new(Entry)
	err = proto.Unmarshal(encoded, entryOut)
	require.NoError(t, err)
	assert.Equal(t, entry, entryOut)
}
