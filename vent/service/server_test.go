// +build integration

package service_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/klyed/hivesmartchain/integration"
	"github.com/klyed/hivesmartchain/integration/rpctest"
	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/vent/config"
	"github.com/klyed/hivesmartchain/vent/service"
	"github.com/klyed/hivesmartchain/vent/sqlsol"
	"github.com/klyed/hivesmartchain/vent/test"
	"github.com/klyed/hivesmartchain/vent/types"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	kern, shutdown := integration.RunNode(t, rpctest.GenesisDoc, rpctest.PrivateAccounts)
	defer shutdown()
	t.Parallel()

	t.Run("Group", func(t *testing.T) {
		t.Run("Run", func(t *testing.T) {
			// run consumer to listen to events
			cfg := config.DefaultVentConfig()

			// create test db
			_, closeDB := test.NewTestDB(t, cfg)
			defer closeDB()

			cfg.SpecFileOrDirs = []string{os.Getenv("GOPATH") + "/src/github.com/klyed/hivesmartchain/vent/test/sqlsol_view.json"}
			cfg.AbiFileOrDirs = []string{os.Getenv("GOPATH") + "/src/github.com/klyed/hivesmartchain/vent/test/EventsTest.abi"}
			cfg.ChainAddress = kern.GRPCListenAddress().String()

			log := logging.NewNoopLogger()
			consumer := service.NewConsumer(cfg, log, make(chan types.EventData))
			projection, err := sqlsol.SpecLoader(cfg.SpecFileOrDirs, sqlsol.None)

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				err := consumer.Run(projection, true)
				require.NoError(t, err)

				wg.Done()
			}()

			time.Sleep(2 * time.Second)

			// setup test server
			server := service.NewServer(cfg, log, consumer)

			httpServer := httptest.NewServer(server)
			defer httpServer.Close()

			// call health endpoint should return OK
			healthURL := fmt.Sprintf("%s/health", httpServer.URL)

			resp, err := http.Get(healthURL)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			// shutdown consumer and wait for its end
			consumer.Shutdown()
			wg.Wait()

			// call health endpoint again should return error
			resp, err = http.Get(healthURL)
			require.NoError(t, err)
			require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		})
	})
}
