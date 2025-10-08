package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"nicolas.galipot.net/hazo/server"
	"nicolas.galipot.net/hazo/server/common"
)

func Serve(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	config := &common.ServerConfig{}

	var addr string
	var debugExplicitlySet bool

	fs.StringVar(&addr, "addr", ":8080", "Server address to listen on (e.g., :8080 or localhost:3000)")
	fs.BoolVar(&config.Debug, "debug", false, "Enable debug mode (auto-enabled for localhost/127.0.0.1)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo serve [options]\n\n")
		fmt.Fprintf(fs.Output(), "Start the web server.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	fs.Visit(func(f *flag.Flag) {
		if f.Name == "debug" {
			debugExplicitlySet = true
		}
	})

	if !debugExplicitlySet {
		addrLower := strings.ToLower(addr)
		if strings.Contains(addrLower, "localhost") || strings.Contains(addrLower, "127.0.0.1") {
			config.Debug = true
		}
	}

	server := server.New(config)
	fmt.Printf("Starting server on %s (debug: %v)\n", addr, config.Debug)
	return http.ListenAndServe(addr, server)
}
