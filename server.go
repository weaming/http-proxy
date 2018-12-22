package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//"github.com/smartystreets/cproxy"
//var explicitForwardProxyHandler = cproxy.Configure()

func SimpleMode(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		//explicitForwardProxyHandler.ServeHTTP(w, r)
		CreateTcpTunnel(w, r.Host, true, nil)
	} else {
		dump := DumpIncomingRequest(r)

		// option 1: exchange normal http
		// ParseDumpAndExchangeReqRes(dump, w, r)

		// option 2: crete tunnel directly
		CreateTcpTunnel(w, getTunnelHost(r), false, &dump)
	}
}

func getTunnelHost(r *http.Request) string{
	host := r.Host
	if !strings.Contains(host, ":") {
		host += ":80"
	}
	return host
}

func ForwardInternet(w http.ResponseWriter, r *http.Request) {
	LogPretty(" ***** ", r.Host)
	if *simple {
		SimpleMode(w, r)
		return
	}

	dump, e := ioutil.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Error: fail read request body"))
		return
	}

	if r.Method == http.MethodConnect {
		CreateTcpTunnel(w, r.Host, true, nil)
	} else {
		ParseDumpAndExchangeReqRes(dump, w, r)
	}
}

func ParseDumpAndExchangeReqRes(dump []byte, w http.ResponseWriter, r *http.Request) {
	//log.Println(string(dump))
	LogPretty(">>> ", strings.SplitN(string(dump), "\n", 2)[0])

	req, e := IncomingRequestToOutgoing(dump, r)
	if e != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Error: fail IncomingRequestToOutgoing"))
		return
	}

	DoRequestAndWriteBack(req, w)
}

func RunServer(listen string) {
	log.Printf("Listening %v\n", listen)
	if err := http.ListenAndServe(listen, http.HandlerFunc(ForwardInternet)); err != nil {
		panic(err)
	}
}
