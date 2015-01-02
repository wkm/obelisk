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

// App is a top-level container which manages the server's execution.
type App struct {
	Config struct {
		StoreDir string
	}

	timedb *storetime.DB // Time series database
	tagdb  *storetag.DB  // Hierarchy database
	kvdb   *storekv.DB   // Key-value database
}

// Start begins
func (app *App) Start() {
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

func (app *App) periodic() {
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker

		// FIXME expose magic number as config
		buffer := make(rinst.MeasurementBuffer, 1000)
		Stats.Measure("obelisk/", &buffer)
	}
}

func (app *App) startTimeStore() (err error) {
	c := storetime.Config{}
	c.DiskStore = filepath.Join(app.Config.StoreDir, "store", "time")
	createPath(c.DiskStore)
	app.timedb, err = storetime.NewDB(c)
	return
}

func (app *App) startTagStore() (err error) {
	c := storetag.Config{}
	c.DiskStore = filepath.Join(app.Config.StoreDir, "store", "tag")
	createPath(c.DiskStore)
	app.tagdb, err = storetag.NewDB(c)
	return
}

func (app *App) startKVStore() (err error) {
	c := storekv.Config{}
	c.DiskStore = filepath.Join(app.Config.StoreDir, "store", "kv")
	createPath(c.DiskStore)
	app.kvdb, err = storekv.NewDB(c)
	return
}

func createPath(path string) (err error) {
	return os.MkdirAll(path, 0777)
}
