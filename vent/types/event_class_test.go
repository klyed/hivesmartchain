package types

import (
	"fmt"
	"testing"

	"github.com/klye-dev/hsc-main/config/source"
)

func TestEventTablesSchema(t *testing.T) {
	schema := ProjectionSpecSchema()
	fmt.Println(source.JSONString(schema))
}
