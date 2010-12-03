// Copyright (c) 2010 Go Fightclub Authors
// Null filter
// As a source: Close without passing any data through the pipeline
// As a link: Pass data straight through without processing
// As a sink: Discard all data

package null

import "afp"

type NullFilter struct {
	ctx *afp.Context
}

func (self *NullFilter) GetType() int {
	return afp.ANY
}

func (self *NullFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx
	parser := flags.FlagParser(args)
	parser.Parse()
}

func (self *NullFilter) Start() {
	header StreamHeader
	ctx := self.ctx
	if (ctx.HeaderSource == nil) { // A source
		ctx.HeaderSink <- StreamHeader{
			Version: 1,
			Channels: 1,
			SampleSize: 0,
			SampleRate: 0,
			ContentLength: 0
		}
		close(ctx.Sink)
	} else if (ctx.HeaderSink == nil) { // A sink
		_ <- ctx.HeaderSource
		for audio := range ctx.Source {
		}
	} else { // A link
		ctx.HeaderSink <- <-ctx.HeaderSource
		for audio := range ctx.Source {
			ctx.Sink <- audio
		}
	}
}
