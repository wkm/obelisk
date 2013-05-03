package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"obelisk/web/util"
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

	resolution, err := strconv.ParseUint(query.Get("resolution"), 10, 64)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	obj["query"] = q
	obj["start"] = start
	obj["resolution"] = resolution
	obj["stop"] = stop

	res, err := QueryTime(q, start/1000, stop/1000)
	if err != nil {
		respondError(rw, err.Error())
	}

	// FIXME extract into a function
	// FIXME should depend on the type of the instrument
	// make into a rate
	var last = math.MaxFloat64
	dps := make([]*util.DataPoint, len(res))
	for i, v := range res {
		val := v.Value
		if last > v.Value {
			val = 0
		} else {
			val = v.Value - last
		}
		last = v.Value

		dsp := util.DataPoint(util.DSPoint{v.Time, val})
		dps[i] = &dsp
	}

	sampled := util.DownSample(uint(resolution), dps)

	points := make([][]interface{}, len(sampled))
	for i, v := range sampled {
		points[i] = make([]interface{}, 3)
		points[i][0] = v.Time * 1000 // ms precision for javascript
		points[i][1] = v.Avg
		points[i][2] = v.Err
	}

	obj["points"] = points

	// response
	rw.Header().Add("Content-Type", "text/json")
	enc := json.NewEncoder(rw)
	enc.Encode(obj)
}
