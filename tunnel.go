package main

import (
	"io"
	"net"
	"net/http"
	"time"
)

// Process CONNECT request and create TCP tunnel
func CreateTcpTunnel(w http.ResponseWriter, host string, writeStatus bool, preData *[]byte) {
	destConn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		LogPretty("  *>> ", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if writeStatus {
		w.WriteHeader(http.StatusOK)
	}
	if preData != nil {
		destConn.Write(*preData)
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()

	io.Copy(destination, source)
}
