package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func timeHandler(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	obj := make(map[string]interface{})

	q := query.Get("query")
	start, err := strconv.ParseUint(query.Get("start"), 10, 64)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	stop, err := strconv.ParseUint(query.Get("stop"), 10, 64)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	obj["query"] = q
	obj["start"] = start
	obj["stop"] = stop

	res, err := QueryTime(q, start/1000, stop/1000)
	if err != nil {
		respondError(rw, err.Error())
	}

	points := make([][]interface{}, len(res))
	for i, v := range res {
		points[i] = make([]interface{}, 2)
		points[i][0] = v.Time * 1000
		points[i][1] = v.Value
	}

	obj["points"] = points

	// response
	rw.Header().Add("Content-Type", "text/json")
	enc := json.NewEncoder(rw)
	enc.Encode(obj)
}
