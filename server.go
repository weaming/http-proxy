package main

import (
	"log"
	"net/http"
	"io/ioutil"
)


func forwardInternet(w http.ResponseWriter, r *http.Request) {
	dump, e := ioutil.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Error: fail read request body"))
		return
	}

	log.Println(string(dump))

	req, e := incomingToOutgoing(dump)
	if e != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Error: fail incomingToOutgoing"))
		return
	}

	doRequest(req, w)
}

func runServer(listen string) {
	http.HandleFunc("/", forwardInternet)
	log.Printf("Listening %v\n", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		panic(err)
	}
}
