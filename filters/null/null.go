// Copyright (c) 2010 Go Fightclub Authors
// The null package defines a set of filters which do nothing, mostly for testing purposes
// NullSource: Close without passing any data through the pipeline
// NullLink: Pass data straight through without processing
// NullSink: Discard all data

package null

import "afp"

//Dummy parent struct, only defines Init
type nullFilter struct {
	ctx *afp.Context
}

func (self *nullFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx
	return nil
}

type NullSource struct {
	nullFilter
}

func NewNullSource() afp.Filter {
	return &NullSource{nullFilter{}}
}

func (self *NullSource) GetType() int {
	return afp.SOURCE
}

func (self *NullSource) Start() {
	self.ctx.HeaderSink <- StreamHeader{
	Version : 1,
	Channels : 1,
	SampleSize : 0,
	SampleRate : 0,
	ContentLength : 0
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
	return afp.SINK
}

func (self *NullSink) Start() {
	_ <- self.ctx.HeaderSource
	for _ := range ctx.Source {
	}
}

type NullLink struct {
	nullFilter
}

func NewNullLink() afp.Filter {
	return &NullLink{nullFilter{}}
}

func (self *NullLink) GetType() int {
	return afp.LINK
}

func (self *NullLink) Start() {
	ctx.HeaderSink <- <-ctx.HeaderSource
	for audio := range ctx.Source {
		ctx.Sink <- audio
	}
}
