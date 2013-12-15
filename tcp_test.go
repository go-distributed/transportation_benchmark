package main

import (
	"io"
	"net"
	"testing"
)

var tcpStarted = false

func BenchmarkTCP(b *testing.B) {
	done := make(chan bool, 10)

	if !tcpStarted {
		go startServer()
	}

	for i := 0; i < b.N; i++ {
		for i := 0; i < 10; i++ {
			go client(done)
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	}
}

func startServer() {
	tcpStarted = true
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go io.Copy(conn, conn)
	}
}

func client(done chan bool) {
	content := make([]byte, 8096)
	for i := 0; i < 8096; i++ {
		content[i] = 32
	}

	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1000; i++ {
		conn.Write(content)
		num, err := conn.Read(content)
		if err != nil || num != 8096 {
			panic(err)
		}
	}

	done <- true
}
