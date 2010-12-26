// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package delay

import (
	"afp"
	"afp/flags"
	"os"
)

type DelayFilter struct {
	context               *afp.Context
	header                afp.StreamHeader
	samplesPerMillisecond int
	delayTimeInMs         int
	delayAttenuation      float32
	extraSamples          int
	bufferSize            int32
	mixBuffer             [][]float32
}

func NewDelayFilter() afp.Filter {
	return &DelayFilter{}
}

func (self *DelayFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *DelayFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.context = ctx

	parser := flags.FlagParser(args)
	var t *int = parser.Int("t", 125, "The delay time in milliseconds")
	var a *float = parser.Float("a", .5, "The wet (delayed) signal amplitude (0 - 1.0)")
	parser.Parse()

	self.delayTimeInMs = *t
	self.delayAttenuation = float32(*a)

	if self.delayTimeInMs <= 0 {
		return os.NewError("Delay time must be greater than zero")
	}

	if self.delayAttenuation < 0 || self.delayAttenuation > 1.0 {
		return os.NewError("Delay signal attenuation must between 0.0 and 1.0")
	}

	return nil
}

func (self *DelayFilter) Start() {
	self.header = <-self.context.HeaderSource
	self.context.HeaderSink <- self.header

	self.samplesPerMillisecond = int(self.header.SampleRate / 1000)
	self.extraSamples = self.delayTimeInMs * self.samplesPerMillisecond

	self.initBuffers()
	self.process()
}

func (self *DelayFilter) process() {
	var (
		t0       int64 = 0
		mbOffset int   = 0
	)

	// loop over all input data
	for audio := range self.context.Source {
		// create a destination buffer
		destBuffer := makeBuffer(self.header.FrameSize, self.header.Channels)

		for t, sample := range audio {
			// mix delayed signal with raw signal
			for c, amplitude := range sample {
				destBuffer[t][c] = amplitude + self.mixBuffer[mbOffset][c]*self.delayAttenuation
			}

			// copy the raw signal into the delay line
			for c, amplitude := range sample {
				self.mixBuffer[mbOffset][c] = amplitude
			}

			// increment the offset into the delay
			mbOffset++
			if mbOffset == self.extraSamples {
				mbOffset = 0
			}
		}

		// send the mixed audio down the pipe
		self.context.Sink <- destBuffer
	}

	// create a destination buffer
	destBuffer := makeBuffer(self.header.FrameSize, self.header.Channels)

	// fill out the rest of the data
	for i := 0; i < self.extraSamples; i++ {
		for c, amplitude := range self.mixBuffer[mbOffset] {
			destBuffer[t0][c] = amplitude * self.delayAttenuation
		}
		t0++
		mbOffset++

		// increment the offset into the delay
		if mbOffset == self.extraSamples {
			mbOffset = 0
		}

		// check to see if we've filled a frame
		if t0 == int64(self.header.FrameSize) {
			// send the mixed audio down the pipe
			self.context.Sink <- destBuffer

			// create a destination buffer
			destBuffer = makeBuffer(self.header.FrameSize, self.header.Channels)
			t0 = 0
		}
	}

	if t0 < int64(self.header.FrameSize) {
		// send the mixed audio down the pipe
		self.context.Sink <- destBuffer
	}
}

// allocate a buffer for samples
func makeBuffer(size int32, channels int8) [][]float32 {
	b := make([][]float32, size)
	for i, _ := range b {
		b[i] = make([]float32, channels)
	}

	return b
}

func (self *DelayFilter) initBuffers() {
	self.mixBuffer = make([][]float32, self.extraSamples)
	for i, _ := range self.mixBuffer {
		self.mixBuffer[i] = make([]float32, self.header.Channels)
	}
}

func (self *DelayFilter) Stop() os.Error {
	return nil
}
