package loggroup

import (
	"fmt"
	"io"
	"sync"
)

// Stream respresent log stream
type Stream interface {
	SetOutput(io.Writer)
	SetPrefix(string)
	SetRowHandler(func(*Row))
	SetCloseHandler(func(Stream))
	GetSequence() []*Row
	Close()

	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

// stream respresent log stream
type stream struct {
	Seq []*Row

	out         io.Writer
	handleRow   func(r *Row)
	handleClose func(s Stream)

	mu       sync.Mutex
	id       int
	level    Level
	isClosed bool
	prefix   string
}

// NewStream return a new Stream
func NewStream(lv Level, out io.Writer) Stream {
	return newStream(lv, out)
}

func newStream(lv Level, out io.Writer) *stream {
	s := &stream{
		Seq: make([]*Row, 0),

		out:         out,
		handleRow:   nil,
		handleClose: nil,

		level:    lv,
		isClosed: false,
	}
	return s
}

// NewNullStream return a useless Stream
func NewNullStream() Stream {
	return NewStream(LevelNone, nil)
}

// SetOutput sets the output destination for the standard logger
func (s *stream) SetOutput(w io.Writer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.out = w
}

// SetPrefix sets the prefix in every log row
func (s *stream) SetPrefix(p string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefix = p
}

// SetRowHandler would be triggered when it append new log Row
func (s *stream) SetRowHandler(handler func(*Row)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handleRow = handler
}

// SetCloseHandler would be triggered when it Closed
func (s *stream) SetCloseHandler(handler func(Stream)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handleClose = handler
}

func (s *stream) GetSequence() []*Row {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]*Row, len(s.Seq))
	copy(res, s.Seq)
	return res
}

// Close closes this stream
func (s *stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed {
		return
	}
	s.isClosed = true
	if s.handleClose != nil {
		s.handleClose(s)
	}
}

func (s *stream) append(lv Level, msg interface{}) {
	if lv < s.level {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed {
		return
	}
	r := newRow(lv, s.prefix, msg)
	s.Seq = append(s.Seq, r)
	if s.out != nil {
		fmt.Fprintln(s.out, r)
	}
	if s.handleRow != nil {
		s.handleRow(r)
	}
}

// Debug record log in debug level
func (s *stream) Debug(msg ...interface{}) {
	s.append(LevelDebug, fmt.Sprint(msg...))
}

// Debugf record log in debug level with format
func (s *stream) Debugf(format string, a ...interface{}) {
	s.append(LevelDebug, fmt.Sprintf(format, a...))
}

// Info record log in info level
func (s *stream) Info(msg ...interface{}) {
	s.append(LevelInfo, fmt.Sprint(msg...))
}

// Infof record log in info level with format
func (s *stream) Infof(format string, a ...interface{}) {
	s.append(LevelInfo, fmt.Sprintf(format, a...))
}

// Warn record log in warning level
func (s *stream) Warn(msg ...interface{}) {
	s.append(LevelWarn, fmt.Sprint(msg...))
}

// Warnf record log in warning level with format
func (s *stream) Warnf(format string, a ...interface{}) {
	s.append(LevelWarn, fmt.Sprintf(format, a...))
}

// Error record log in error level
func (s *stream) Error(msg ...interface{}) {
	s.append(LevelError, fmt.Sprint(msg...))
}

// Errorf record log in error level with format
func (s *stream) Errorf(format string, a ...interface{}) {
	s.append(LevelError, fmt.Errorf(format, a...))
}
