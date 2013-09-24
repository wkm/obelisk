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

//
func (app *ServerApp) ReceiveStats(worker string, buffer rinst.MeasurementBuffer) error {
	for {
		measure, ok := <-buffer
		if !ok {
			return nil
		}

		err := app.RecordMeasurement(worker, measure)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *ServerApp) ReceiveStatsBuffered(worker string, buffer []rinst.Measurement) error {
	for _, measure := range buffer {
		err := app.RecordMeasurement(worker, measure)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *ServerApp) RecordMeasurement(worker string, measure rinst.Measurement) error {
	_, err := app.tagdb.Store.NewTag("tag/" + measure.Name)
	if err != nil {
		return err
	}

	id, err := app.tagdb.Store.NewTag("worker/" + worker + "/" + measure.Name)
	flt, err := strconv.ParseFloat(measure.Value, 64)
	if err != nil {
		log.Printf("invalid measurement %s in %s", err.Error(), measure)
		return nil
	}

	app.timedb.Store.Insert(id, measure.Time, flt)
	return nil
}

func (app *ServerApp) DeclareSchema(worker string, buffer rinst.SchemaBuffer) error {
	log.Printf("receiving schema from %s", worker)
	for {
		schema, ok := <-buffer
		if !ok {
			return nil
		}

		err := app.RecordSchema(worker, schema)
		if err != nil {
			return err
		}
	}
}

// Declare a schema through a single slice buffer (instead of a channel)
func (app *ServerApp) DeclareSchemaBuffered(worker string, buffer []rinst.Schema) error {
	for _, schema := range buffer {
		err := app.RecordSchema(worker, schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *ServerApp) RecordSchema(worker string, schema rinst.Schema) error {
	err := app.kvdb.Store.SetGob("worker/"+worker+"/"+schema.Name, schema)
	return err
}
