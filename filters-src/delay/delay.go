// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package delay

import (
	"libav"
	"afp/types"
	"flags"
)

type DelayFilter struct {
	ctx *types.Context
	samplesPerSecond int = 44100
	samplesPerMillisecond int = samplesPerSecond / 1000
	delayTimeInMs int = 150
	channels int16 = 2
	bytesPerSample int16 = 2
}

func (self *DelayFilter) GetType() int {
	return types.PIPE_LINK
}

func Usage(args []string) os.Error {
	msg = fmt.Sprintf(os.Stderr, "Usage of %s:\n"
		"%s <delay per millisecond>\n", args[0], args[0])
	return os.NewError(msg)


func (self *DelayFilter) Init(ctx *types.Context, args []string) os.Error {
	self.ctx = ctx
	if len(vars) != 2 {
		return Usage(args)
	}
	self.samplesPerSecond, err := strconv.Atoui(args[1])
	if err != nil {
		return os.NewError("No delay specified")
	}
	return nil
}


func (self *DelayFilter) Start() {
	// Some of this will be moved to FileSource
	self.header <-self.context.HeaderSource
	var avcontext libav.AVDecodeContext
	var currBuffer = 0
	buffer := make([][]float32, ctx.channels)
}
