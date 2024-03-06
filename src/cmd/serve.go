package cmd

import "nicolas.galipot.net/hazo/server"

func Serve(addr string) error {
	return server.Serve(addr)
}
