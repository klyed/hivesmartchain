package query

import (
	"testing"

	"github.com/klyed/hivesmartchain/logging/structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder(t *testing.T) {
	qb := NewBuilder()
	qry, err := qb.Query()
	require.NoError(t, err)
	assert.Equal(t, "", qry.String())

	qb = qb.AndGreaterThanOrEqual("foo.size", 45)
	qry, err = qb.Query()
	require.NoError(t, err)
	assert.Equal(t, "foo.size >= 45", qry.String())

	qb = qb.AndEquals("bar.name", "marmot")
	qry, err = qb.Query()
	require.NoError(t, err)
	assert.Equal(t, "foo.size >= 45 AND bar.name = 'marmot'", qry.String())

	assert.True(t, qry.Matches(makeTagMap("foo.size", 80, "bar.name", "marmot")))
	assert.False(t, qry.Matches(makeTagMap("foo.size", 8, "bar.name", "marmot")))
	assert.False(t, qry.Matches(makeTagMap("foo.size", 80, "bar.name", "marot")))

	qb = qb.AndContains("bar.desc", "hsc")
	qry, err = qb.Query()
	require.NoError(t, err)
	assert.Equal(t, "foo.size >= 45 AND bar.name = 'marmot' AND bar.desc CONTAINS 'hsc'", qry.String())

	assert.True(t, qry.Matches(makeTagMap("foo.size", 80, "bar.name", "marmot", "bar.desc", "lives in a hsc")))
	assert.False(t, qry.Matches(makeTagMap("foo.size", 80, "bar.name", "marmot", "bar.desc", "lives in a shoe")))

	qb = NewBuilder().AndEquals("foo", "bar")
	qb = qb.And(NewBuilder().AndGreaterThanOrEqual("frogs", 4))
	qry, err = qb.Query()
	require.NoError(t, err)
	assert.Equal(t, "foo = 'bar' AND frogs >= 4", qry.String())

	qb = qb.Not()
	qry, err = qb.Query()
	require.NoError(t, err)
	assert.Equal(t, "NOT (foo = 'bar' AND frogs >= 4)", qry.String())
}

func makeTagMap(keyvals ...interface{}) TagMap {
	tmap := make(TagMap)
	for i := 0; i < len(keyvals); i += 2 {
		tmap[keyvals[i].(string)] = structure.Stringify(keyvals[i+1])
	}
	return tmap
}
