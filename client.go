package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ForwardServer(w http.ResponseWriter, r *http.Request) {
	LogPretty(" ***** ", r.Host)
	dump := DumpIncomingRequest(r)
	LogPretty(">>> ", strings.SplitN(string(dump), "\n", 2)[0])

	if r.Method == http.MethodConnect {
		CreateTcpTunnel(w, *server, false, &dump)
		return
	}

	req, e := http.NewRequest("POST", *server, bytes.NewReader(dump))
	if e != nil {
		w.WriteHeader(500)
		log.Fatal(w.Write([]byte("500 Internal Error: fail on NewRequest")))
		return
	}

	AddForwardedHeaders(req, r)
	DoRequestAndWriteBack(req, w)
}

func AddForwardedHeaders(req, originReq *http.Request) {
	fBy := "package-via-http"
	fFor := originReq.Header.Get("User-Agent")
	fHost := originReq.Host
	fProto := originReq.Proto
	//fSchema := originReq.URL.Scheme
	fSchema := "https"
	ip := originReq.RemoteAddr
	LogPretty("  >>> ", originReq.Header)

	forward := fmt.Sprintf("by=%v; for=%v; host=%v; proto=%v", fBy, fFor, fHost, fProto)
	req.Header.Set("Forwarded", forward)
	req.Header.Set("X-Forwarded-By", fBy)
	req.Header.Set("X-Forwarded-For", fFor)
	req.Header.Set("X-Forwarded-Host", fHost)
	req.Header.Set("X-Forwarded-Proto", fProto)
	req.Header.Set("X-Forwarded-Schema", fSchema)
	req.Header.Set("X-Real-IP", ip)
}

func RunClient(listen, server string) {
	log.Printf("Listening %v\n", listen)
	outerHandler := http.HandlerFunc(ForwardServer)
	log.Fatal(http.ListenAndServe(listen, outerHandler))
}
