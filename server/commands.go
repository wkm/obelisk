package server

// Declare associates an identifier with the given paths
func (app *App) Declare(id string, paths ...string) (err error) {
	statDeclare.Incr()
	return
	// for _, path := range paths {
	// 	actid, err := app.tagdb.Tag(id, path)
	// 	if id != actid {
	// 		return errors.New("identifier is not unique within path")
	// 	}
	// }
}

// Schema stores metadata on the a metric's structure.
func (app *App) Schema(id, op, kind, unit, desc string) (err error) {
	statSchema.Incr()
	return
}

// Record stores a single measurement of a metric.
func (app *App) Record(id, metric, time, value string) (err error) {
	statRecord.Incr()
	return
}
