package main

import (
	"bytes"
	"circuit/sys/zanchorfs"
	"encoding/gob"
	"errors"
	"fmt"
)

// try to decode an arbitrary blog into some interface, and string print it
func getAsGob(blob string) (string, error) {
	if len(blob) < 1 {
		return "", errors.New("no content")
	}

	var thing struct{}
	err := gob.NewDecoder(bytes.NewBufferString(blob)).Decode(&thing)
	if err != nil {
		return "", err
	}

	// if thing == nil {
	// 	return "", errors.New("not a gob thing")
	// }

	return fmt.Sprintf("%#v", thing), nil
}

// try to derive an anchor file from a blob
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
