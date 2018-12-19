package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"bytes"
	"bufio"
	"io/ioutil"
	"log"
)

var Client = http.DefaultClient

func dumpRequest(req *http.Request) []byte {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	return dump
}

func incomingToOutgoing(c []byte) (*http.Request, error) {
	r, e := http.ReadRequest(bufio.NewReader(bytes.NewReader(c)))
	if e != nil {
		return r, e
	}
	r.RequestURI = ""
	return r, nil
}

func doRequest(req *http.Request, w http.ResponseWriter) {
	res, e := Client.Do(req)
	if e != nil {
		log.Printf("%+v\n", req)
		w.WriteHeader(502)
		lg := fmt.Sprintf("502 Bad Gateway: fail doing request: %v\n", e)
		log.Println(lg)
		w.Write([]byte(lg))
		return
	}

	// read body
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		w.WriteHeader(502)
		lg := "502 Bad Gateway: fail read response body"
		log.Println(lg)
		w.Write([]byte(lg))
		return
	}

	// status
	w.WriteHeader(res.StatusCode)
	logPretty("<<< status: ", res.StatusCode)

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
	logPretty("<<< headers: ", w.Header())

	// body
	w.Write(body)
}

func logPretty(prefix string, v interface{}) {
	log.Printf("%v%+v\n", prefix, v)
}
