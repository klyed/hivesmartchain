// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package loggers

import (
	"encoding"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/klyed/hivesmartchain/logging/structure"
	"github.com/tmthrgd/go-hex"
)

// Logger that implements some formatting conventions for hsc and hsc-client
// This is intended for applying consistent value formatting before the final 'output' logger;
// we should avoid prematurely formatting values here if it is useful to let the output logger
// decide how it wants to display values. Ideal candidates for 'early' formatting here are types that
// we control and generic output loggers are unlikely to know about.
type hscFormatLogger struct {
	sync.Mutex
	logger  log.Logger
	options opt
}

type opt byte

func (o opt) enabled(q opt) bool {
	return o&q > 0
}

const (
	DefaultOptions opt = iota
	StringifyValues
)

func NewHiveSmartChainFormatLogger(logger log.Logger, options ...opt) *hscFormatLogger {
	bfl := &hscFormatLogger{logger: logger}
	for _, option := range options {
		bfl.options |= option
	}
	return bfl
}

var _ log.Logger = &hscFormatLogger{}

func (bfl *hscFormatLogger) Log(keyvals ...interface{}) error {
	if bfl.logger == nil {
		return nil
	}
	keyvals, err := structure.MapKeyValues(keyvals,
		func(key interface{}, value interface{}) (interface{}, interface{}) {
			switch v := value.(type) {
			case string, json.Marshaler, encoding.TextMarshaler:
			case time.Time:
				value = v.Format(time.RFC3339Nano)
			case fmt.Stringer:
				value = v.String()
			case []byte:
				value = hex.EncodeUpperToString(v)
			}
			if bfl.options.enabled(StringifyValues) {
				value = structure.Stringify(value)
			}
			return structure.Stringify(key), value
		})
	if err != nil {
		return err
	}
	bfl.Lock()
	defer bfl.Unlock()
	return bfl.logger.Log(keyvals...)
}
