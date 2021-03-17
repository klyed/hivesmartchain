package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBurrow(t *testing.T) {
	var outputCount int
	out := &output{
		PrintfFunc: func(format string, args ...interface{}) {
			outputCount++
		},
		LogfFunc: func(format string, args ...interface{}) {
			outputCount++
		},
		FatalfFunc: func(format string, args ...interface{}) {
			t.Fatalf("fatalf called by Hive Smart Chain cmd: %s", fmt.Sprintf(format, args...))
		},
	}
	app := hsc(out)
	// Basic smoke test for cli config
	err := app.Run([]string{"hsc", "--version"})
	assert.NoError(t, err)
	err = app.Run([]string{"hsc", "spec", "--name-prefix", "foo", "-f1"})
	assert.NoError(t, err)
	err = app.Run([]string{"hsc", "configure"})
	assert.NoError(t, err)
	err = app.Run([]string{"hsc", "start", "--help"})
	assert.NoError(t, err)
	assert.True(t, outputCount > 0)
}
