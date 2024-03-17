package cmd

import (
	"net/http"

	"nicolas.galipot.net/hazo/server"
)

func Serve(args []string) error {
	addr := args[0]
	server := server.New()
	return http.ListenAndServe(addr, server)
}
