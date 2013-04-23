package main

import (
	"bytes"
	"circuit/sys/zanchorfs"
	"encoding/gob"
	"errors"
)

func getAnchorFile(blob string) (*zanchorfs.ZFile, error) {
	if len(blob) < 1 {
		return nil, errors.New("no content")
	}

	afile := &zanchorfs.ZFile{}
	err := gob.NewDecoder(bytes.NewBufferString(blob)).Decode(afile)
	if err != nil {
		return nil, err
	}

	return afile, nil
}
