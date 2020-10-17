package loggroup

import (
	"fmt"
	"time"
)

// Row defines log row information
type Row struct {
	Time   time.Time
	Level  Level
	Prefix string
	Msg    interface{}
}

// newRow return a new log row
func newRow(lv Level, prefix string, msg interface{}) *Row {
	return &Row{
		Time:   time.Now(),
		Level:  lv,
		Prefix: prefix,
		Msg:    msg,
	}
}

func (r *Row) String() string {
	return fmt.Sprintf("%-32s [%-5s] %s %s",
		r.Time.Format(time.RFC3339Nano),
		r.Level,
		r.Prefix,
		r.Msg)
}
