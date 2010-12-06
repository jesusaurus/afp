// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.
//
// The null package defines a set of filters which do nothing, mostly for testing purposes
// NullSource: Close without passing any data through the pipeline
// NullLink: Pass data straight through without processing
// NullSink: Discard all data

package null

import (
	"afp"
	"os"
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
}

func NewNullSource() afp.Filter {
	return &NullSource{nullFilter{}}
}

func (self *NullSource) GetType() int {
	return afp.PIPE_SOURCE
}

func (self *NullSource) Start() {
	self.ctx.HeaderSink <- afp.StreamHeader{
	Version : 1,
	Channels : 1,
	SampleSize : 0,
	SampleRate : 0,
	ContentLength : 0,
	}
	close(self.ctx.Sink)
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