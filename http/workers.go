package main

import (
	"bytes"
	"circuit/sys/zanchorfs"
	// "circuit/use/anchorfs"
	"encoding/gob"
	"log"
)

type WorkerStatus uint

var (
	StatusOk      = 0
	StatusError   = 1
	StatusUnknown = 2
)

type Worker struct {
	id, host string
	status   WorkerStatus
}

func getAnchorFile(blob string) (*zanchorfs.ZFile, error) {
	afile := &zanchorfs.ZFile{}
	err := gob.NewDecoder(bytes.NewBufferString(blob)).Decode(afile)
	if err != nil {
		log.Printf("error parsing anchor file: %s", err)
		return nil, err
	}

	log.Printf("anchor file: %s", afile.Addr)
	return afile, nil
}
