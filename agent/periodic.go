package agent

// this is the core of the agent; executed on a regularly scheduled period
func Periodic() {
	statMeasurements.Incr()
}
