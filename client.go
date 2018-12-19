package main

import (
	"net/http"
	"log"
	"bytes"
	"strings"
	"github.com/smartystreets/cproxy"
)

var explicitForwardProxyHandler = cproxy.Configure()

func forwardServer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		//w.WriteHeader(200)
		explicitForwardProxyHandler.ServeHTTP(w, r)
		return
	}

	dump := dumpRequest(r)
	logPretty(">>> ", strings.SplitN(string(dump), "\n", 2)[0])

	req, e := http.NewRequest("POST", *server, bytes.NewReader(dump))
	//dumpReq := dumpRequest(req)
	//log.Println(string(dumpReq))
	if e != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Error: fail on NewRequest"))
		return
	}

	doRequest(req, w)
}

func runClient(listen, server string) {
	log.Printf("Listening %v\n", listen)
	outerHandler := http.HandlerFunc(forwardServer)
	http.ListenAndServe(listen, outerHandler)
}
