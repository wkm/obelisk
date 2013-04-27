package agent

// this is the core of the agent; executed on a regularly scheduled period
func Periodic() {
	time.Sleep(5 * time.Second)
	rlog.Log.Printf("it is now %s\n", time.Now().Format(time.Kitchen))
}
