// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package structure

import (
	"testing"

	"github.com/klyed/hivesmartchain/util/slice"

	"github.com/stretchr/testify/assert"
)

func TestValuesAndContext(t *testing.T) {
	keyvals := []interface{}{"hello", 1, "dog", 2, "fish", 3, "fork", 5}
	vals, ctx := ValuesAndContext(keyvals, "hello", "fish")
	assert.Equal(t, map[string]interface{}{"hello": 1, "fish": 3}, vals)
	assert.Equal(t, []interface{}{"dog", 2, "fork", 5}, ctx)
}

func TestKeyValuesMap(t *testing.T) {
	keyvals := []interface{}{
		[][]interface{}{{2}}, 3,
		"hello", 1,
		"fish", 3,
		"dog", 2,
		"fork", 5,
	}
	vals := KeyValuesMap(keyvals)
	assert.Equal(t, map[string]interface{}{
		"[[2]]": 3,
		"hello": 1,
		"fish":  3,
		"dog":   2,
		"fork":  5,
	}, vals)
}

func TestVectorise(t *testing.T) {
	kvs := []interface{}{
		"scope", "lawnmower",
		"hub", "budub",
		"occupation", "fish brewer",
		"scope", "hose pipe",
		"flub", "dub",
		"scope", "rake",
		"flub", "brub",
	}

	kvsVector := Vectorise(kvs, "occupation", "scope")
	// Vectorise scope
	assert.Equal(t, []interface{}{
		"scope", Vector{"lawnmower", "hose pipe", "rake"},
		"hub", "budub",
		"occupation", "fish brewer",
		"flub", Vector{"dub", "brub"},
	},
		kvsVector)
}

func TestVector_String(t *testing.T) {
	vec := Vector{"one", "two", "grue"}
	assert.Equal(t, "[one two grue]", vec.String())
}

func TestRemoveKeys(t *testing.T) {
	// Remove multiple of same key
	assert.Equal(t, []interface{}{"Fish", 9},
		RemoveKeys([]interface{}{"Foo", "Bar", "Fish", 9, "Foo", "Baz", "odd-key"},
			"Foo"))

	// Remove multiple different keys
	assert.Equal(t, []interface{}{"Fish", 9},
		RemoveKeys([]interface{}{"Foo", "Bar", "Fish", 9, "Foo", "Baz", "Bar", 89},
			"Foo", "Bar"))

	// Remove nothing but supply keys
	assert.Equal(t, []interface{}{"Foo", "Bar", "Fish", 9},
		RemoveKeys([]interface{}{"Foo", "Bar", "Fish", 9},
			"A", "B", "C"))

	// Remove nothing since no keys supplied
	assert.Equal(t, []interface{}{"Foo", "Bar", "Fish", 9},
		RemoveKeys([]interface{}{"Foo", "Bar", "Fish", 9}))
}

func TestDelete(t *testing.T) {
	assert.Equal(t, []interface{}{1, 2, 4, 5}, Delete([]interface{}{1, 2, 3, 4, 5}, 2, 1))
}

func TestCopyPrepend(t *testing.T) {
	assert.Equal(t, []interface{}{"three", 4, 1, "two"},
		slice.CopyPrepend([]interface{}{1, "two"}, "three", 4))
	assert.Equal(t, []interface{}{}, slice.CopyPrepend(nil))
	assert.Equal(t, []interface{}{1}, slice.CopyPrepend(nil, 1))
	assert.Equal(t, []interface{}{1}, slice.CopyPrepend([]interface{}{1}))
}
