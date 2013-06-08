package util

import (
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"obelisk/lib/errors"
	"obelisk/lib/rlog"
	"obelisk/server"
)

var log = rlog.LogConfig.Logger("obelisk-reporter")

// through zookeeper discover an obelisk server and connect to it
func DiscoverObeliskServer() (circuit.X, error) {
	log.Printf("discovering obelisk server")
	nodes, err := anchorfs.OpenDir("/obelisk-server")
	if err != nil {
		log.Printf("could not find /obelisk-server %s", err.Error())
		return nil, err
	}

	_, workers, err := nodes.Files()
	if err != nil {
		log.Printf("could not list workers %s", err.Error())
		return nil, err
	}

	if len(workers) < 1 {
		err = errors.N("could not find any obelisk-server workers")
		log.Printf(err.Error())
		return nil, err
	}

	log.Printf("found obelisk-server workers %#v", workers)
	var xServer circuit.X
	for id, file := range workers {
		xServer, err = circuit.TryDial(file.Owner(), server.ServiceName)
		if err != nil {
			log.Printf("  error dialing %v:%v with %s", id, file, err.Error())
		} else {
			// stop at the first usable connection
			break
		}
	}

	if xServer == nil {
		err = errors.N("could not find a responsive obelisk-server worker")
		return nil, err
	}

	log.Printf("connected to %#v", xServer)
	return xServer, nil
}
