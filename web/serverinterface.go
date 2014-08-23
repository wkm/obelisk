package main

import (
	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/storetime"
)

func ChildrenTags(node ...string) (tags []string, err error) {
	return
	// retrn := xServer.Call("ChildrenTags", filepath.Join(node...))
	// var children []string
	// var err error
	// if retrn[0] != nil {
	// 	children = retrn[0].([]string)
	// }
	// if retrn[1] != nil {
	// 	err = retrn[1].(error)
	// }

	// return children, err
}

func QueryTime(node string, start, stop uint64) (points []storetime.Point, err error) {
	return
	// retrn := xServer.Call("QueryTime", node, start, stop)
	// var values []storetime.Point
	// var err error
	// if retrn[0] != nil {
	// 	values = retrn[0].([]storetime.Point)
	// }
	// if retrn[1] != nil {
	// 	err = retrn[1].(error)
	// }
	// return values, err
}

func GetMetricInfo(node string) (schema rinst.InstrumentSchema, err error) {
	return
	// retrn := xServer.Call("GetMetricInfo", node)
	// var info rinst.InstrumentSchema
	// var err error
	// if retrn[0] != nil {
	// 	info = retrn[0].(rinst.InstrumentSchema)
	// }
	// if retrn[1] != nil {
	// 	err = retrn[1].(error)
	// }
	// return info, err
}
