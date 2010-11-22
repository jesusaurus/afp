package afp

import (
	"os"
	"log"
)

//Constants to specify the type of a given filter
const (
	PIPE_SOURCE = iota
	PIPE_SINK
	PIPE_LINK
)

type StreamHeader struct {
	Version       int
	Channels      int
	SampleSize    int
	SampleRate    int
	ContentLength int64
}

type Context struct {
	HeaderSource <-chan StreamHeader
	HeaderSink   chan<- StreamHeader
	Source       <-chan [][]float32
	Sink         chan<- [][]float32

	Verbose   bool
	Err, Info *log.Logger
}

type Filter interface {
	GetType() int
	Init(*Context, []string) os.Error
	Start()
}
