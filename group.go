package loggroup

import (
	"io"
	"os"
	"sync"
)

// Group defines log group
type Group interface {
	Stream
	NewStream() Stream
	GetLevel() Level
}

// group defines log group
type group struct {
	Stream
	mu sync.Mutex

	level        Level
	idCounter    int
	out          io.Writer
	handleStream func(s Stream)
	handleRow    func(r *Row)
}

// New return a new group with default configuration
func New(lv Level) Group {
	return NewGroupWithHandle(lv, os.Stderr, nil, nil)
}

// NewNullGroup return a useless Group
func NewNullGroup() Group {
	return NewGroupWithHandle(LevelNone, nil, nil, nil)
}

// NewGroup return a new group
func NewGroup(lv Level, w io.Writer) Group {
	return NewGroupWithHandle(lv, w, nil, nil)
}

// NewGroupWithHandle return a new group
func NewGroupWithHandle(lv Level, w io.Writer, handleStream func(Stream), handleRow func(*Row)) Group {
	g := &group{
		level:        lv,
		out:          w,
		handleStream: handleStream,
		handleRow:    handleRow,
	}
	g.Stream = g.NewStream()
	return g
}

// NewStream return a new Stream belong to this Group
func (g *group) NewStream() Stream {
	s := newStream(g.level, g.out)
	s.SetRowHandler(g.handleRow)
	s.SetCloseHandler(g.handleStream)
	g.idCounter++
	s.id = g.idCounter
	return s
}

// GetLevel return log level
func (g *group) GetLevel() Level {
	return g.level
}
