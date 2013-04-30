package main

import (
	"bytes"
	"circuit/sys/zanchorfs"
	"encoding/gob"
	"fmt"
	"obelisk/lib/errors"
)

// try to decode an arbitrary blog into some interface, and string print it
func getAsGob(blob string) (string, error) {
	if len(blob) < 1 {
		return "", errors.N("no content")
	}

	var thing struct{}
	// err := gob.NewDecoder(bytes.NewBufferString(blob)).Decode(&thing)
	if err != nil {
		return "", err
	}

	// if thing == nil {
	// 	return "", errors.N("not a gob thing")
	// }

	return fmt.Sprintf("%#v", thing), nil
}

// try to derive an anchor file from a blob
func getAnchorFile(blob string) (*zanchorfs.ZFile, error) {
	if len(blob) < 1 {
		return nil, errors.N("no content")
	}

	afile := &zanchorfs.ZFile{}
	err := gob.NewDecoder(bytes.NewBufferString(blob)).Decode(afile)
	if err != nil {
		return nil, errors.W(err)
	}

	return afile, nil
}
