package rlog

// the null log doesn't do any actual logging
type NullLog struct{}

func (l NullLog) Printf(format string, obj ...interface{}) {}
func (l NullLog) Sync()                                    {}
func (l NullLog) Close()                                   {}
