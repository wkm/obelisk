package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/wkm/obelisk/web/util"
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

	rate, err := strconv.ParseBool(query.Get("rate"))
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	obj["query"] = q
	obj["start"] = start
	obj["resolution"] = resolution
	obj["stop"] = stop
	obj["rate"] = rate

	// get all data points from the
	res, err := QueryTime(q, start/1000, stop/1000)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	// FIXME extract into a function
	// FIXME we should be able to the rate thing in-place
	var last = util.DSPoint{0, math.MaxFloat64}
	dps := make([]*util.DataPoint, len(res))
	for i, v := range res {
		val := v.Value
		if rate {
			if last.Value() > v.Value {
				val = 0
			} else {
				val = (v.Value - last.V) / float64(v.Time-last.T)
			}

			last.T = v.Time
			last.V = v.Value
		}

		dsp := util.DataPoint(util.DSPoint{v.Time, val})
		dps[i] = &dsp
	}

	sampled := util.DownSample(start/1000, stop/1000, uint(resolution), dps)
	points := make([][]interface{}, len(sampled))
	for i, v := range sampled {
		points[i] = make([]interface{}, 3)
		points[i][0] = v.Time * 1000 // ms precision for javascript

		if !math.IsNaN(v.Avg) {
			points[i][1] = v.Avg
		} else {
			points[i][1] = nil
		}

		if !math.IsNaN(v.Err) {
			points[i][2] = v.Err
		}
	}

	obj["points"] = points

	rw.Header().Add("Content-Type", "text/json")
	enc := json.NewEncoder(rw)
	err = enc.Encode(obj)
	if err != nil {
		respondError(rw, err.Error())
		return
	}
}
