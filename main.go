package main

import (
	"flag"
)

var listen = flag.String("listen", ":10008", "linsten on")
var server = flag.String("server", "", "server host:port")
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
