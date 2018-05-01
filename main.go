package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type webRequest struct {
	r      *http.Request
	w      http.ResponseWriter
	doneCh chan struct{}
}

var (
	requestCh = make(chan *webRequest)
)

var (
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
)

func init() {
	http.DefaultClient = &http.Client{Transport: &transport}
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doneCh := make(chan struct{})

		request := webRequest{
			r:      r,
			w:      w,
			doneCh: doneCh,
		}

		requestCh <- &request

		<-doneCh
	})

	go processRequests()

	go http.ListenAndServe(":2000", nil)

	println("Load balancer stared, pres <Enter> to exit.")

	fmt.Scanln()
}
