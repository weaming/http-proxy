package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"strings"
)

func ForwardInternet(w http.ResponseWriter, r *http.Request) {
	dump, e := ioutil.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Error: fail read request body"))
		return
	}

	if r.Method == http.MethodConnect {
		ServerProcessTcpTunnel(w, r)
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
	http.HandleFunc("/", ForwardInternet)
	log.Printf("Listening %v\n", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		panic(err)
	}
}
