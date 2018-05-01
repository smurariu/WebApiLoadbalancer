package main

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	appservers   = []string{}
	currentIndex = 0
	client       = http.Client{Transport: &transport, Timeout: 10 * time.Second}
)

func processRequests() {
	for {
		select {
		case request := <-requestCh:
			println("got a request")
			if len(appservers) == 0 {
				request.w.WriteHeader(http.StatusInternalServerError)
				request.w.Write([]byte("No appservers were found."))
				request.doneCh <- struct{}{}
				continue
			}

			currentIndex++
			if currentIndex == len(appservers) {
				currentIndex = 0
			}

			appserverURL := appservers[currentIndex]
			go processRequest(appserverURL, request)

		}
	}
}

func processRequest(appserverURL string, request *webRequest) {
	hostURL, _ := url.Parse(request.r.URL.String())
	hostURL.Scheme = "http"
	hostURL.Host = appserverURL
	println(appserverURL)
	println(hostURL.String())

	//new up a request and add the headers
	req, _ := http.NewRequest(request.r.Method, hostURL.String(), request.r.Body)
	for k, v := range request.r.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		req.Header.Add(k, values)
	}

	resp, err := client.Do(req)
	if err != nil {
		request.w.WriteHeader(http.StatusInternalServerError)
		request.doneCh <- struct{}{}
		return
	}

	for k, v := range resp.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}

		request.w.Header().Add(k, values)
	}

	io.Copy(request.w, resp.Body)

	request.doneCh <- struct{}{}
}