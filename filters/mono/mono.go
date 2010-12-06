// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.
//
// Mono filter
// Combine all channels into a single channel
// Uses the average of all channels

package mono

import (
	"afp"
	"os"
)

type MonoFilter struct {
	ctx *afp.Context
}

func (self *MonoFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *MonoFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx
	if len(os.Args) != 1 {
		return os.NewError("mono link takes 0 arguments")
	}
	return nil
}

func (self *MonoFilter) Start() {
	var header afp.StreamHeader
	ctx := self.ctx
	header = <-ctx.HeaderSource
	// Unpack header
	channels := header.Channels

	header.Channels = 1
	// TODO: Is this math correct? Is this guaranteed to be accurate?
	header.ContentLength = header.ContentLength / int64(channels)
	ctx.HeaderSink <- header

	if channels == 1 { // Already mono, don't manipulate
		for frame := range ctx.Source {
			ctx.Sink <- frame
		}
	} else {
		for frame := range ctx.Source {
			ctx.Sink <- mergeChannels(frame, channels)
		}
	}
}


func mergeChannels(frame [][]float32, channels int8) [][]float32 {
	var monoValue float32
	for sample, sampleValues := range frame {
		monoValue = 0.0
		for _, channelValue := range sampleValues {
			monoValue += channelValue / float32(channels)
		}

		frame[sample] = []float32{monoValue}
	}
	return frame
}
