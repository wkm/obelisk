package server

import (
	"circuit/sys/transport"
	"circuit/use/circuit"
	"obelisk/lib/errors"
	"obelisk/lib/rinst"
	"obelisk/lib/rinst/service"
	"obelisk/lib/storekv"
	"obelisk/lib/storetag"
	"obelisk/lib/storetime"
	"path/filepath"
	"time"
)

type ServerApp struct {
	timedb *storetime.DB
	tagdb  *storetag.DB
	kvdb   *storekv.DB
}

const ObeliskDirectory = "/var/run/obelisk"

func (app *ServerApp) Main() {
	log.Printf("starting obelisk-server")
	app.startTimeStore()
	app.startTagStore()
	app.startKVStore()

	service.Expose(Stats)

	// FIXME all this probably should use the flush utilities
	buffer := make(rinst.SchemaBuffer, 1000)
	go func() {
		Stats.Schema("", buffer)
		close(buffer)
	}()

	app.RegisterWorker(circuit.WorkerAddr())
	app.DeclareSchema(circuit.WorkerAddr().WorkerID().String(), buffer)

	go app.periodic()
}

func (app *ServerApp) periodic() {
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker

		// FIXME expose magic number as config
		buffer := make(rinst.MeasurementBuffer, 1000)
		go func() {
			Stats.Measure("", buffer)
			close(buffer)
		}()

		// report our own stats
		worker := circuit.WorkerAddr().WorkerID().String()
		app.ReceiveStats(worker, buffer)
	}
}

func (app *ServerApp) startTimeStore() {
	var err error
	c := storetime.NewConfig()
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "time")
	app.timedb, err = storetime.NewDB(c)
	if err != nil {
		log.Printf("error starting time store %s", err.Error())
	}
}

func (app *ServerApp) startTagStore() {
	var err error
	c := storetag.NewConfig()
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "tag")
	app.tagdb, err = storetag.NewDB(c)
	if err != nil {
		log.Printf("error starting tag store %s", err.Error())
	}
}

func (app *ServerApp) startKVStore() {
	var err error
	c := storekv.NewConfig()
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "kv")
	app.kvdb, err = storekv.NewDB(c)
	if err != nil {
		log.Printf("error starting keyvalue store%s", err.Error())
	}
}

func (app *ServerApp) ChildrenTags(node string) ([]string, error) {
	return app.tagdb.Store.Children(node)
}

func (app *ServerApp) QueryTime(node string, start, stop uint64) ([]storetime.Point, error) {
	id, err := app.tagdb.Store.Id(node)
	if err != nil {
		return nil, errors.W(err)
	}

	return app.timedb.Store.Query(id, start, stop)
}

func (app *ServerApp) GetWorkerAddr(worker string) (circuit.Addr, error) {
	var addr transport.Addr
	err := app.kvdb.Store.GetGob("worker/"+worker+"/addr", &addr)
	return &addr, errors.W(err)
}

func (app *ServerApp) GetMetricInfo(node string) (rinst.Schema, error) {
	var schema rinst.Schema
	err := app.kvdb.Store.GetGob(node, &schema)
	return schema, errors.W(err)
}
