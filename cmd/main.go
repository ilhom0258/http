package main

import (
	"net"
	"os"

	"github.com/ilhom0258/server/pkg/server"
)

func main() {
	if err := execute(server.Host, server.Port); err != nil {
		os.Exit(1)
	}
}

func execute(host, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", srv.RouteHandler("Welcome to our web-site"))
	srv.Register("/about", srv.RouteHandler("About Golang Academy"))
	return srv.Start()
}
