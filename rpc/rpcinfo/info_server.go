// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package rpcinfo

import (
	"net"
	"net/http"

	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/logging/structure"
	"github.com/klyed/hivesmartchain/rpc"
	"github.com/klyed/hivesmartchain/rpc/lib/server"
)

func StartServer(service *rpc.Service, pattern string, listener net.Listener, logger *logging.Logger) (*http.Server, error) {
	logger = logger.With(structure.ComponentKey, "RPC_Info")
	routes := GetRoutes(service)
	mux := http.NewServeMux()
	server.RegisterRPCFuncs(mux, routes, logger)
	srv, err := server.StartHTTPServer(listener, mux, logger)
	if err != nil {
		return nil, err
	}
	return srv, nil
}
