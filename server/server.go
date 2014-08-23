package server

import (
	"path/filepath"
	"time"

	"github.com/wkm/obelisk/lib/errors"
	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/storekv"
	"github.com/wkm/obelisk/lib/storetag"
	"github.com/wkm/obelisk/lib/storetime"
)

type ServerApp struct {
	timedb *storetime.DB // time series database
	tagdb  *storetag.DB  // hierarchy database
	kvdb   *storekv.DB   // key-value database
}

const ObeliskDirectory = "/var/run/obelisk"

func (app *ServerApp) Main() {
	log.Printf("starting obelisk-server")
	app.startTimeStore()
	app.startTagStore()
	app.startKVStore()

	// FIXME all this probably should use the flush utilities
	buffer := make(rinst.SchemaBuffer, 1000)
	Stats.Schema("", &buffer)

	go app.periodic()
}

func (app *ServerApp) periodic() {
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker

		// FIXME expose magic number as config
		buffer := make(rinst.MeasurementBuffer, 1000)
		Stats.Measure("", &buffer)
	}
}

func (app *ServerApp) startTimeStore() {
	var err error
	c := storetime.Config{}
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "time")
	app.timedb, err = storetime.NewDB(c)
	if err != nil {
		log.Printf("error starting time store %s", err.Error())
	}
}

func (app *ServerApp) startTagStore() {
	var err error
	c := storetag.Config{}
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "tag")
	app.tagdb, err = storetag.NewDB(c)
	if err != nil {
		log.Printf("error starting tag store %s", err.Error())
	}
}

func (app *ServerApp) startKVStore() {
	var err error
	c := storekv.Config{}
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "kv")
	app.kvdb, err = storekv.NewDB(c)
	if err != nil {
		log.Printf("error starting keyvalue store%s", err.Error())
	}
}

func (app *ServerApp) ChildrenTags(node string) ([]string, error) {
	return app.tagdb.Children(node)
}

// query a time series between start and stop
func (app *ServerApp) QueryTime(node string, start, stop uint64) ([]storetime.Point, error) {
	id, err := app.tagdb.Id(node)
	if err != nil {
		return nil, errors.W(err)
	}

	return app.timedb.Query(id, start, stop)
}

func (app *ServerApp) GetMetricInfo(node string) (rinst.InstrumentSchema, error) {
	var schema rinst.InstrumentSchema
	err := app.kvdb.GetGob(node, &schema)
	return schema, errors.W(err)
}
