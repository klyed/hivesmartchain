// +build integration,!ethereum

package service_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/klyed/hivesmartchain/integration"
	"github.com/klyed/hivesmartchain/integration/rpctest"

	"github.com/stretchr/testify/assert"

	"github.com/klyed/hivesmartchain/vent/types"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/klyed/hivesmartchain/vent/test"
)

func TestPostgresConsumer(t *testing.T) {
	privateAccounts := rpctest.PrivateAccounts
	kern, shutdown := integration.RunNode(t, rpctest.GenesisDoc, privateAccounts)
	defer shutdown()
	inputAddress := privateAccounts[0].GetAddress()
	grpcAddress := kern.GRPCListenAddress().String()
	tcli := test.NewHiveSmartChainTransactClient(t, grpcAddress)

	t.Parallel()
	time.Sleep(2 * time.Second)

	t.Run("Group", func(t *testing.T) {
		t.Run("PostgresConsumer", func(t *testing.T) {
			testConsumer(t, kern.Blockchain.ChainID(), test.PostgresVentConfig(grpcAddress), tcli, inputAddress)
		})

		t.Run("PostgresInvalidUTF8", func(t *testing.T) {
			testInvalidUTF8(t, test.PostgresVentConfig(grpcAddress), tcli, inputAddress)
		})

		t.Run("PostgresDeleteEvent", func(t *testing.T) {
			testDeleteEvent(t, kern.Blockchain.ChainID(), test.PostgresVentConfig(grpcAddress), tcli, inputAddress)
		})

		t.Run("PostgresResume", func(t *testing.T) {
			testResume(t, test.PostgresVentConfig(grpcAddress))
		})

		t.Run("PostgresTriggers", func(t *testing.T) {
			tCli := test.NewHiveSmartChainTransactClient(t, kern.GRPCListenAddress().String())
			create := test.CreateContract(t, tCli, inputAddress)

			// generate events
			name := "TestTriggerEvent"
			description := "Trigger it!"
			txe := test.CallAddEvent(t, tCli, inputAddress, create.Receipt.ContractAddress, name, description)

			cfg := test.PostgresVentConfig(grpcAddress)
			// create test db
			_, closeDB := test.NewTestDB(t, cfg)
			defer closeDB()

			// Create a postgres notification listener
			listener := pq.NewListener(cfg.DBURL, time.Second, time.Second*20, func(event pq.ListenerEventType, err error) {
				require.NoError(t, err)
			})

			// These are defined in sqlsol_view.json
			err := listener.Listen("meta")
			require.NoError(t, err)

			err = listener.Listen("keyed_meta")
			require.NoError(t, err)

			err = listener.Listen(types.BlockHeightLabel)
			require.NoError(t, err)

			type payload struct {
				Height uint64 `json:"_height"`
			}

			heightCh := make(chan uint64)
			notifications := make(map[string]string)
			go func() {
				for n := range listener.Notify {
					notifications[n.Channel] = n.Extra
					if n.Channel == types.BlockHeightLabel {
						pl := new(payload)
						err := json.Unmarshal([]byte(n.Extra), pl)
						if err != nil {
							panic(err)
						}
						if pl.Height >= txe.Height {
							heightCh <- pl.Height
							return
						}
					}
				}
			}()
			resolveSpec(cfg, testViewSpec)
			runConsumer(t, cfg)

			// Give events a chance
			const timeout = 3 * time.Second
			select {
			case <-time.After(timeout):
				t.Fatalf("timed out after %v waiting for notification", timeout)
			case height := <-heightCh:
				// Assert we get expected returns
				t.Logf("latest height: %d, txe height: %d", height, txe.Height)
				assert.True(t, height >= txe.Height)
			}
			assert.Equal(t, `{"_action" : "INSERT", "testdescription" : "5472696767657220697421000000000000000000000000000000000000000000", "testname" : "TestTriggerEvent"}`,
				notifications["meta"])
			assert.Equal(t, `{"_action" : "INSERT", "testdescription" : "5472696767657220697421000000000000000000000000000000000000000000", "testkey" : "\\x544553545f4556454e5453000000000000000000000000000000000000000000", "testname" : "TestTriggerEvent"}`,
				notifications["keyed_meta"])
		})
	})
}
