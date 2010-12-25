// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.
//
// The null package defines a set of filters which do nothing, mostly for testing purposes
// NullSource: Output silence
// NullLink: Pass data straight through without processing
// NullSink: Discard all data

package null

import (
	"afp"
	"os"
	"afp/flags"
)
//Dummy parent struct, only defines Init/Stop
type nullFilter struct {
	ctx *afp.Context
}

func (self *nullFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	return nil
}

func (self *nullFilter) Stop() os.Error {
	return nil
}

type NullSource struct {
	nullFilter
	time, samplerate, framesize, channels int
}

func NewNullSource() afp.Filter {
	return &NullSource{nullFilter : nullFilter{}}
}

func (self *NullSource) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	parser := flags.FlagParser(args)
	parser.IntVar(&self.time, "t", 10, "Time in seconds of silence to output")
	parser.IntVar(&self.samplerate, "s", 44100, "Sample rate to output.")
	parser.IntVar(&self.framesize, "f", 256, "Frame size to output.")
	parser.IntVar(&self.channels, "c", 2, "Number of channels in output signal")
	parser.Parse()

	if self.time < 0 {
		os.NewError("Time must be greater than 0.")
	}

	if self.samplerate < 1 {//This should probably be higher
		os.NewError("Sample rate must be at least 1.")
	}

	if self.framesize < 1 {
		os.NewError("Frame size must be at least 1.")
	}

	if self.channels < 1 {	
		os.NewError("Channels must be at least 1.")
	}

	return nil
}

func (self *NullSource) GetType() int {
	return afp.PIPE_SOURCE
}

func (self *NullSource) Start() {
	length := int64(self.time) * int64(self.samplerate) //Overflow paranoia

	self.ctx.HeaderSink <- afp.StreamHeader{
	Version:       1,
	Channels:      int8(self.channels),
	SampleSize:    4,
	SampleRate:    int32(self.samplerate),
	FrameSize:     int32(self.framesize),
	ContentLength: length,
	}

	frames := length / int64(self.framesize)
	if length % int64(self.framesize) != 0 {
		frames++
	}
		
	for i := int64(0); i < frames; i++ {
		raw := make([]float32, self.framesize * self.channels)
		f := make([][]float32, self.framesize)

		for i, s, e := 0, 0, self.channels; i < self.framesize; i++ {
			f[i] = raw[s:e]
			s = e
			e += self.channels
		}

		self.ctx.Sink <- f
	}
}

type NullSink struct {
	nullFilter
}

func NewNullSink() afp.Filter {
	return &NullSink{nullFilter{}}
}

func (self *NullSink) GetType() int {
	return afp.PIPE_SINK
}

func (self *NullSink) Start() {
	<-self.ctx.HeaderSource
	for _ = range self.ctx.Source {
		//Do nothing
	}
}

type NullLink struct {
	nullFilter
}

func NewNullLink() afp.Filter {
	return &NullLink{nullFilter{}}
}

func (self *NullLink) GetType() int {
	return afp.PIPE_LINK
}

func (self *NullLink) Start() {
	self.ctx.HeaderSink <- <-self.ctx.HeaderSource
	for audio := range self.ctx.Source {
		self.ctx.Sink <- audio
	}
}
