package types

import (
	"fmt"
	"testing"

	"github.com/klye-dev/hivesmartchain/config/source"
)

func TestEventTablesSchema(t *testing.T) {
	schema := ProjectionSpecSchema()
	fmt.Println(source.JSONString(schema))
}
