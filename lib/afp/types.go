// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

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
	Version       int8
	Channels      int8
	SampleSize    int8
	SampleRate    int32
	FrameSize     int32
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
	Stop() os.Error
}
