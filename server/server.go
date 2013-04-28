package server

import (
	"path/filepath"
	// "circuit/use/circuit"
	"log"
	"obelisk/lib/storekv"
	"obelisk/lib/storetag"
	"obelisk/lib/storetime"
)

const ServiceName = "obelisk-server"

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
}

func (app *ServerApp) startTimeStore() {
	var err error
	c := storetime.NewConfig()
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "time")
	app.timedb, err = storetime.NewDB(c)
	if err != nil {
		log.Fatalf("error starting time store %s", err.Error())
	}
}

func (app *ServerApp) startTagStore() {
	var err error
	c := storetag.NewConfig()
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "tag")
	app.tagdb, err = storetag.NewDB(c)
	if err != nil {
		log.Fatalf("error starting tag store %s", err.Error())
	}
}

func (app *ServerApp) startKVStore() {
	var err error
	c := storekv.NewConfig()
	c.DiskStore = filepath.Join(ObeliskDirectory, "store", "kv")
	app.kvdb, err = storekv.NewDB(c)
	if err != nil {
		log.Fatalf("error starting keyvalue store%s", err.Error())
	}
}
