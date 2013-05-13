package server

import (
	"circuit/use/circuit"
	"encoding/gob"
	"obelisk/lib/rinst"
	"obelisk/lib/rlog"
	"strconv"
)

var log = rlog.LogConfig.Logger("obelisk-server")

func init() {
	gob.Register(rinst.Schema{})
}

func (app *ServerApp) RegisterWorker(addr circuit.Addr) {
	worker := addr.WorkerID().String()
	app.kvdb.Store.SetGob("worker/"+worker+"/addr", addr)
	app.tagdb.Store.NewTag("host/" + addr.Host() + "/workers/" + worker)
}

func (app *ServerApp) ReceiveStats(worker string, buffer rinst.MeasurementBuffer) error {
	log.Printf("receiving stats from %s", worker)
	for {
		measure, ok := <-buffer
		if !ok {
			return nil
		}

		_, err := app.tagdb.Store.NewTag("tag/" + measure.Name)
		if err != nil {
			return err
		}

		id, err := app.tagdb.Store.NewTag("worker/" + worker + "/" + measure.Name)
		flt, err := strconv.ParseFloat(measure.Value, 64)
		if err != nil {
			log.Printf("invalid measurement %s in %s", err.Error(), measure)
			continue
		}

		app.timedb.Store.Insert(id, measure.Time, flt)
	}
}

func (app *ServerApp) DeclareSchema(worker string, buffer rinst.SchemaBuffer) error {
	log.Printf("receiving schema from %s", worker)
	for {
		schema, ok := <-buffer
		if !ok {
			return nil
		}

		err := app.kvdb.Store.SetGob("worker/"+worker+"/"+schema.Name, schema)
		if err != nil {
			return err
		}
	}
}
