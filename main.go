package main

import (
	"flag"
)

var listen = flag.String("listen", ":10008", "linsten on")
var server = flag.String("server", "", "server host:port")

func init() {
	flag.Parse()
}

func main() {
	if *server == "" {
		runServer(*listen)
	} else {
		runClient(*listen, *server)
	}
}
