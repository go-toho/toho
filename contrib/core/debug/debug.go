package debug

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // Enables pprof endpoint.
)

type Config struct {
	Enabled bool   `default:"true"`
	Addr    string `default:":6060"`
}

func NewHTTPServer(config Config) (*http.Server, error) {
	host, port, err := net.SplitHostPort(config.Addr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve address: %v", err)
	}

	if host == "" {
		host = "localhost"
	}

	server := &http.Server{Addr: fmt.Sprintf("%s:%s", host, port)}

	return server, nil
}
