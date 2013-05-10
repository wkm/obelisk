package server

import (
	"encoding/gob"
	"obelisk/lib/rinst"
	"obelisk/lib/rlog"
	"strconv"
)

var log = rlog.LogConfig.Logger("obelisk-server")

func init() {
	gob.Register(rinst.Schema{})
}

func (app *ServerApp) ReceiveStats(hostname string, buffer rinst.MeasurementBuffer) error {
	log.Printf("receiving stats from %s", hostname)
	for {
		measure, ok := <-buffer
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

func (app *ServerApp) DeclareSchema(hostname string, buffer rinst.SchemaBuffer) error {
	log.Printf("receiving schema from %s", hostname)
	for {
		schema, ok := <-buffer
		if !ok {
			return nil
		}

		err := app.kvdb.Store.SetGob("host/"+hostname+"/"+schema.Name, schema)
		if err != nil {
			return err
		}
	}
}
