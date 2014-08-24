package server

import (
	"os"
	"path/filepath"
	"time"

	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/storekv"
	"github.com/wkm/obelisk/lib/storetag"
	"github.com/wkm/obelisk/lib/storetime"
)

type ServerApp struct {
	StoreDir string

	timedb *storetime.DB // time series database
	tagdb  *storetag.DB  // hierarchy database
	kvdb   *storekv.DB   // key-value database
}

func (app *ServerApp) Start() {
	log.Printf("Starting")

	if err := app.startTimeStore(); err != nil {
		log.Printf("Couldn't start time store %s", err.Error())
		os.Exit(-1)
	}

	if err := app.startTagStore(); err != nil {
		log.Printf("Couldn't start tag store %s", err.Error())
		os.Exit(-1)
	}

	if err := app.startKVStore(); err != nil {
		log.Printf("Couldn't start key-value store %s", err.Error())
		os.Exit(-1)
	}

	// Server self reports
	buffer := make(rinst.SchemaBuffer, 0)
	Stats.Schema("obelisk/", &buffer)
	go app.periodic()
}

func (app *ServerApp) periodic() {
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker

		// FIXME expose magic number as config
		buffer := make(rinst.MeasurementBuffer, 1000)
		Stats.Measure("obelisk/", &buffer)
	}
}

func (app *ServerApp) startTimeStore() (err error) {
	c := storetime.Config{}
	c.DiskStore = filepath.Join(app.StoreDir, "store", "time")
	app.timedb, err = storetime.NewDB(c)
	return
}

func (app *ServerApp) startTagStore() (err error) {
	c := storetag.Config{}
	c.DiskStore = filepath.Join(app.StoreDir, "store", "tag")
	app.tagdb, err = storetag.NewDB(c)
	return
}

func (app *ServerApp) startKVStore() (err error) {
	c := storekv.Config{}
	c.DiskStore = filepath.Join(app.StoreDir, "store", "kv")
	app.kvdb, err = storekv.NewDB(c)
	return
}
