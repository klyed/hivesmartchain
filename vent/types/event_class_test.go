package types

import (
	"fmt"
	"testing"

	"github.com/klyed/hivesmartchain/config/source"
)

func TestEventTablesSchema(t *testing.T) {
	schema := ProjectionSpecSchema()
	fmt.Println(source.JSONString(schema))
}
