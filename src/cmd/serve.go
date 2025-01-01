package cmd

import (
	"net/http"

	"nicolas.galipot.net/hazo/server"
	"nicolas.galipot.net/hazo/server/common"
)

func Serve(args []string, config *common.ServerConfig) error {
	addr := args[0]
	server := server.New(config)
	return http.ListenAndServe(addr, server)
}
