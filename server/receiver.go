package server

import (
	"log"
	"obelisk/lib/rinst"
	"strconv"
)

func (app *ServerApp) ReceiveStats(hostname string, measurements rinst.MeasurementBuffer) error {
	log.Printf("receiving stats from %s", hostname)
	for {
		measure, ok := <-measurements
		if !ok {
			return nil
		}

		_, err := app.tagdb.Store.NewTag("tag/" + measure.Name)
		if err != nil {
			return err
		}

		id, err := app.tagdb.Store.NewTag("host/" + hostname + "/" + measure.Name)
		flt, err := strconv.ParseFloat(measure.Value, 64)
		if err != nil {
			log.Printf("invalid measurement %s in %s", err.Error(), measure)
			continue
		}

		app.timedb.Store.Insert(id, measure.Time, flt)
	}
}
