package server

import (
	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/rlog"
)

var log = rlog.LogConfig.Logger("obelisk-server")

// ... replace with a proper endpoint
func (app *ServerApp) ReceiveStats(worker string, r rinst.MeasurementBuffer) error {
	for _, measure := range r {
		err := app.RecordMeasurement(worker, measure)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *ServerApp) RecordMeasurement(worker string, measure rinst.InstrumentMeasurement) error {
	_, err := app.tagdb.NewTag("tag/" + measure.Name)
	if err != nil {
		return err
	}

	id, err := app.tagdb.NewTag("worker/" + worker + "/" + measure.Name)
	if measure.IntValue != 0 {
		app.timedb.Insert(id, uint64(measure.Time), float64(measure.IntValue))
	} else {
		app.timedb.Insert(id, uint64(measure.Time), measure.FloatValue)
	}

	return nil
}

// Declare a schema through a single slice buffer (instead of a channel)
func (app *ServerApp) DeclareSchema(worker string, buffer rinst.SchemaBuffer) error {
	for _, schema := range buffer {
		err := app.RecordSchema(worker, schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *ServerApp) RecordSchema(worker string, schema rinst.InstrumentSchema) error {
	err := app.kvdb.SetGob("worker/"+worker+"/"+schema.Name, schema)
	return err
}
