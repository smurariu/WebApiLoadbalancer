package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type webRequest struct {
	r      *http.Request
	w      http.ResponseWriter
	doneCh chan struct{}
}

var (
	requestCh    = make(chan *webRequest)
	registerCh   = make(chan string)
	unregisterCh = make(chan string)
	heartbeat    = time.Tick(5 * time.Second)
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

	go http.ListenAndServe(":2002", new(appserverHandler))

	println("Load balancer stared, pres <Enter> to exit.")

	fmt.Scanln()
}

type appserverHandler struct{}

func (h *appserverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := ""

	if strings.HasPrefix(r.RemoteAddr, "[::1]") {
		ip = "localhost"
	} else {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	port := r.URL.Query().Get("port")
	switch r.URL.Path {
	case "/register":
		registerCh <- ip + ":" + port

	case "/unregister":
		unregisterCh <- ip + ":" + port
	}
}
