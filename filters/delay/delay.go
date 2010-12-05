// Copyright (c) 2010 Go Fightclub Authors

package delay

import (
	"libav"
	"afp/types"
	"flags"
)

type DelayFilter struct {
	ctx *afp.Context
	header *afp.StreamHeader
	samplesPerSecond int = 44100
	samplesPerMillisecond int = samplesPerSecond / 1000
	delayTimeInMs int = 150
	extraSamples int
	bufferSize int
	channels int16 = 2
	bytesPerSample int16 = 2
	buffers [][][]float
	frameCopy [][]float
}

func (self *DelayFilter) GetType() int {
	return types.PIPE_LINK
}

func Usage(args []string) os.Error {
	msg = fmt.Sprintf(os.Stderr, "Usage of %s:\n"
		"%s <delay per millisecond>\n", args[0], args[0])
	return os.NewError(msg)
}

func (self *DelayFilter) Init(ctx *types.Context, args []string) os.Error {
	self.ctx = ctx
	if len(vars) != 2 {
		return Usage(args)
	}
	
	self.delayTimeInMs, err := strconv.Atoui(args[1])
	if err != nil {
		return os.NewError("No delay specified")
	}
	
	self.extraSamples = self.delayTimeInMs * self.samplesPerMillisecond;

	return nil
}

func (self *DelayFilter) Start() {
	// Some of this will be moved to FileSource
	self.header <-self.context.HeaderSource
	self.initBuffers()
	
	self.process()
}

func (self *DelayFilter) process() {
	var currBuffer := 0
	
	
}

func (self *DelayFilter) initBuffers() {
	self.bufferSize = ctx.FrameSize

	self.buffers := make([][][]float32, 2)
	for _,buffer := range buffers {
		buffer := make([][]float32, ctx.FrameSize)

		for _,sample := range buffer {
			sample := make([]float32, ctx.Channels)
		}
	}
	
	self.frameCopy := make([][]float32, ctx.FrameSize)
	for _,sample := range frameCopy {
		sample := make([]float32, ctx.Channels)
	}
}
