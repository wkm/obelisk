package main

import (
	"flag"
	"net"

	"github.com/wkm/obelisk/lib/resp"
	"github.com/wkm/obelisk/lib/rlog"
	"github.com/wkm/obelisk/server"
)

var log = rlog.LogConfig.Logger("obelisk-server")

var (
	storeDir = flag.String("data", "/tmp/obelisk", "directory to place data stores")
	address  = flag.String("address", ":6666", "address to listen on for connections")
)

func main() {
	s := new(server.App)
	s.Config.StoreDir = *storeDir
	s.Start()

	// Start listening
	ln, err := net.Listen("tcp", *address)
	defer ln.Close()
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Listening for connection on %s", *address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting conn: %s", err.Error())
			continue
		}

		log.Printf("Connection from %s", conn.RemoteAddr())
		go resp.Listen(s, conn, conn)
	}
}
