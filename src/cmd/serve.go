package cmd

import "nicolas.galipot.net/hazo/server"

func Serve(args []string) error {
	addr := args[0]
	return server.Serve(addr)
}
