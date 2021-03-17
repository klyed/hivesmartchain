package sqlsol_test

import (
	"os"
	"path"
	"testing"

	"github.com/KLYE-Dev/HSC-MAIN/vent/sqlsol"
	"github.com/KLYE-Dev/HSC-MAIN/vent/types"
	"github.com/stretchr/testify/require"
)

var tables = types.DefaultSQLTableNames

func TestSpecLoader(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)
	specFile := []string{path.Join(dir, "../test/sqlsol_view.json")}
	t.Run("successfully add block and transaction tables to event structures", func(t *testing.T) {
		projection, err := sqlsol.SpecLoader(specFile, sqlsol.BlockTx)
		require.NoError(t, err)

		require.Equal(t, 4, len(projection.Tables))

		require.Equal(t, tables.Block, projection.Tables[tables.Block].Name)

		require.Equal(t, columns.Height,
			projection.Tables[tables.Block].GetColumn(columns.Height).Name)

		require.Equal(t, tables.Tx, projection.Tables[tables.Tx].Name)

		require.Equal(t, columns.TxHash,
			projection.Tables[tables.Tx].GetColumn(columns.TxHash).Name)
	})
}
