package main

import (
	"flag"
)

var listen = flag.String("listen", ":8123", "listen host:port")
var server = flag.String("server", "", "server host:port, run as client if given")
var simple = flag.Bool("simple", false, "serve http proxy without client")

func init() {
	flag.Parse()
}

func main() {
	if *server == "" {
		RunServer(*listen)
	} else {
		RunClient(*listen, *server)
	}
}
