package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"strings"
)

func ForwardInternet(w http.ResponseWriter, r *http.Request) {
	LogPretty(" ***** ", r.Host)
	if *simple {
		if r.Method == http.MethodConnect {
			//explicitForwardProxyHandler.ServeHTTP(w, r)
			ProcessTcpTunnel(w, r.Host, true, nil)
		} else {
			dump := DumpIncomingRequest(r)
			ParseDumpAndExchangeReqRes(dump, w, r)
		}
		return
	}

	dump, e := ioutil.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Error: fail read request body"))
		return
	}

	if r.Method == http.MethodConnect {
		ProcessTcpTunnel(w, r.Host, true, nil)
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
