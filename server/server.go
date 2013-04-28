package server

import (
	"log"
	"obelisk/lib/rinst"
	"obelisk/lib/storekv"
	"obelisk/lib/storetag"
	"obelisk/lib/storetime"
	"os"
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

	go app.periodic()
}

func (app *ServerApp) periodic() {
	ticker := time.Tick(1 * time.Minute)
	for {
		<-ticker

		// FIXME expose magic number as config
		buffer := make(rinst.MeasurementBuffer, 1000)
		go func() {
			Stats.Measure("", buffer)
			close(buffer)
		}()

		host, err := os.Hostname()
		if err != nil {
			log.Printf("could not get hostname %s", err.Error())
			continue
		}

		app.ReceiveStats(host, buffer)
	}
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
