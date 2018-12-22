package main

import (
	"flag"
)

var listen = flag.String("listen", ":8123", "listen host:port")
var server = flag.String("server", "", "server host:port, run as client if given")
var simple = flag.Bool("simple", false, "serve http proxy without client")
var key = flag.String("key", "", "key to encrypt data between client and server")

// proxy chain
var httpProxy = flag.String("http", "", "http proxy host:port")
var httpsProxy = flag.String("https", "", "https proxy host:port")
var socksV4Proxy = flag.String("socks4", "", "socks4 proxy host:port")
var socksV5Proxy = flag.String("socks5", "", "socks5 proxy host:port")

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
