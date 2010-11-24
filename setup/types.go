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
	ANY
)

type StreamHeader struct {
	HeaderLength  int32
	Version       int8
	Channels      int8
	SampleSize    int8
	SampleRate    int32
	FrameSize     int32
	ContentLength int64
	OtherLength   int32
	Other         []byte
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
