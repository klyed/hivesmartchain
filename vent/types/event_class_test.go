package types

import (
	"fmt"
	"testing"

	"github.com/KLYE-Dev/HSC-MAIN/config/source"
)

func TestEventTablesSchema(t *testing.T) {
	schema := ProjectionSpecSchema()
	fmt.Println(source.JSONString(schema))
}
