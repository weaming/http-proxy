package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"bytes"
	"bufio"
	"io/ioutil"
	"log"
	"strings"
)

var Client = http.DefaultClient

func DumpIncomingRequest(req *http.Request) []byte {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	return dump
}

func IncomingRequestToOutgoing(c []byte, originReq *http.Request) (*http.Request, error) {
	req, e := http.ReadRequest(bufio.NewReader(bytes.NewReader(c)))
	if e != nil {
		return req, e
	}

	//LogPretty("  >>> ", originReq.Header)

	// uri
	req.RequestURI = ""
	// schema
	schema := originReq.Header.Get("X-Forwarded-Schema")
	if schema == "" {
		// simple mode
		schema = "https"
	}
	req.URL.Scheme = schema
	// proxy connection
	req.Header.Del("Proxy-Connection")

	return req, nil
}

func DoRequestAndWriteBack(req *http.Request, w http.ResponseWriter) {
	LogPretty("  *>> ", req.Header)
	res, e := Client.Do(req)
	if e != nil {
		//LogPretty("  *>> ", req)
		w.WriteHeader(502)
		lg := fmt.Sprintf("* 502 Bad Gateway: fail doing request: %v\n", e)
		log.Println(lg)
		w.Write([]byte(lg))
		return
	}

	// read body
	length := res.ContentLength
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		w.WriteHeader(502)
		lg := fmt.Sprintf("* 502 Bad Gateway: fail reading response body: %v\n", e)
		log.Println(lg)
		w.Write([]byte(lg))
		return
	}

	// status
	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Length", fmt.Sprintf("%v", length))
	LogPretty("<<< status: ", res.StatusCode)

	// headers
	for k, vs := range res.Header {
		for i, v := range vs {
			if i == 0 {
				w.Header().Set(k, v)
			} else {
				w.Header().Add(k, v)
		}
		}
	}
	LogPretty("<<< headers: ", w.Header())

	// body
	w.Write(body)
}

func LogPretty(prefix string, v interface{}) {
	log.Printf("%v%+v\n", prefix, v)
}

func correctTcpHost(host string) string {
	if !strings.Contains(host, ":") {
		host += ":80"
	}
	return host
}
func correctHttpHost(host string) string {
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	return host
}
