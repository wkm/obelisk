package rlog

// NullLog is a black hole for logging information.
type NullLog struct{}

// Printf is a NOP with a NullLog.
func (l NullLog) Printf(format string, obj ...interface{}) {}

// Sync is a NOP with a NullLog.
func (l NullLog) Sync() {}

// Close is a NOP with a NullLog.
func (l NullLog) Close() {}
