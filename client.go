package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ForwardServer(w http.ResponseWriter, r *http.Request) {
	LogPretty("***** CRH ", r.Host)
	dump := DumpIncomingRequest(r)
	LogPretty("  *>> ", strings.SplitN(string(dump), "\n", 2)[0])

	if r.Method == http.MethodConnect {
		// pass by the CONNECT to server
		CreateTcpTunnel(w, *server, false, &dump)
		return
	}

	if req, e := http.NewRequest("POST", correctHttpHost(*server), bytes.NewReader(dump)); e != nil {
		w.WriteHeader(500)
		log.Fatal(w.Write([]byte("500 Internal Error: fail on NewRequest")))
	} else {
		DoRequestAndWriteBack(AddForwardedHeaders(req, r), w)
	}
}

func AddForwardedHeaders(req, originReq *http.Request) *http.Request {
	fBy := "package-via-http"
	fFor := originReq.Header.Get("User-Agent")
	fHost := originReq.Host
	fProto := originReq.Proto
	//fSchema := originReq.URL.Scheme
	fSchema := "https"
	ip := originReq.RemoteAddr
	//LogPretty("  *>> ", originReq.Header)

	forward := fmt.Sprintf("by=%v; for=%v; host=%v; proto=%v", fBy, fFor, fHost, fProto)
	req.Header.Set("Forwarded", forward)
	req.Header.Set("X-Forwarded-By", fBy)
	req.Header.Set("X-Forwarded-For", fFor)
	req.Header.Set("X-Forwarded-Host", fHost)
	req.Header.Set("X-Forwarded-Proto", fProto)
	req.Header.Set("X-Forwarded-Schema", fSchema)
	req.Header.Set("X-Real-IP", ip)

	return req
}

func RunClient(listen, server string) {
	log.Printf("Listening %v\n", listen)
	outerHandler := http.HandlerFunc(ForwardServer)
	log.Fatal(http.ListenAndServe(listen, outerHandler))
}
