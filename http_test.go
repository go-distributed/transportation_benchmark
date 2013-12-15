package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var httpStarted = false

func BenchmarkHTTP(b *testing.B) {
	done := make(chan bool, 10)

	if !httpStarted {
		go startHTTPServer()
	}

	for i := 0; i < b.N; i++ {
		for i := 0; i < 10; i++ {
			go httpClient(done)
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	}
}
func startHTTPServer() {
	httpStarted = true
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	http.ListenAndServe(":6666", nil)
}

func httpClient(done chan bool) {
	for i := 0; i < 1000; i++ {
		resp, err := http.Get("http://127.0.0.1:6666/")
		if err != nil {
			panic(err)
		}
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	done <- true
}
