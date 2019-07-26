package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"

	"github.com/atotto/clipboard"
)

var (
	addr = flag.String("addr", "127.0.0.1:5010", "Address to listen on")
)

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal("failed to listen:  %v", err)
	}
	log.Printf("started server on %v\n", ln.Addr())

	err = processConnections(ln)
	log.Fatal("process ended: %v", err)
}

func processConnections(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			// Failed to accept, is it temporary?
			type isTemporary interface {
				Temporary() bool
			}

			if err, ok := err.(isTemporary); ok && err.Temporary() {
				log.Printf("Ignoring temporary accept error: %v", err)
				continue
			}

			return err
		}

		go worker(conn)
	}

}

func worker(c net.Conn) {
	bs, err := ioutil.ReadAll(c)
	log.Printf("Got connection with (truncated) text: %.20s\n", bs)
	clipboard.WriteAll(string(bs))
	if err != nil {
		log.Printf("  failed to read successfully: %v", err)
	}
}
