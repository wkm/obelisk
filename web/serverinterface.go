package main

import (
	"circuit/use/circuit"
	"obelisk/lib/rinst"
	"obelisk/lib/storetime"
	"path/filepath"
)

var xServer circuit.X

func ChildrenTags(node ...string) ([]string, error) {
	retrn := xServer.Call("ChildrenTags", filepath.Join(node...))
	var children []string
	var err error
	if retrn[0] != nil {
		children = retrn[0].([]string)
	}
	if retrn[1] != nil {
		err = retrn[1].(error)
	}

	return children, err
}

func QueryTime(node string, start, stop uint64) ([]storetime.Point, error) {
	retrn := xServer.Call("QueryTime", node, start, stop)
	var values []storetime.Point
	var err error
	if retrn[0] != nil {
		values = retrn[0].([]storetime.Point)
	}
	if retrn[1] != nil {
		err = retrn[1].(error)
	}
	return values, err
}

func GetMetricInfo(node string) (rinst.Schema, error) {
	retrn := xServer.Call("GetMetricInfo", node)
	var info rinst.Schema
	var err error
	if retrn[0] != nil {
		info = retrn[0].(rinst.Schema)
	}
	if retrn[1] != nil {
		err = retrn[1].(error)
	}
	return info, err
}
