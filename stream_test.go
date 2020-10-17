package loggroup

import (
	"fmt"
	"testing"
)

func TestStream(t *testing.T) {
	s := NewStream(LevelDebug, nil)
	s.Error(fmt.Errorf("This is Error"))
	s.Info("This is Info")
	s.Warn("This is Warn")
	s.Debug("This is Debug")

	s.Close()
	s.Debug("This is useless")

	if len(s.GetSequence()) != 4 {
		t.Fatal("stream len error")
	}
}
