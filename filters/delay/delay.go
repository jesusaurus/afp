// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package delay

import (
	"afp"
	"flags"
	"os"
)

type DelayFilter struct {
	context *afp.Context
	header afp.StreamHeader
	samplesPerSecond int
	samplesPerMillisecond int
	delayTimeInMs int
	extraSamples int
	bufferSize int32
	channels int16
	bytesPerSample int16
	buffers [][][]float32
	frameCopy [][]float32
}

func (self *DelayFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *DelayFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.context = ctx

	parser := flags.FlagParser(args)
	var t *int = parser.Int("t", 125, "The delay time in milliseconds")
/*	var w *int = parser.Int("w", 40, "The wet (delayed) signal ratio: 0 (dry) to 100 (wet)")*/
	parser.Parse()
	
	self.delayTimeInMs = *t	
	self.extraSamples = self.delayTimeInMs * self.samplesPerMillisecond;

	return nil
}

func (self *DelayFilter) Start() {
	self.header = <-self.context.HeaderSource
	self.context.HeaderSink <- self.header
	
	self.initBuffers()
	self.process()
}

func (self *DelayFilter) process() {
	var (
		t int64 = 0
		t1 int64 = 0
		d float32 = 0.75
		w float32 = 0.25
	)

	buffer := 0
	destBuffer := self.buffers[buffer][:] 
	
	for audio := range(self.context.Source) {
		self.frameCopy = audio
		for _,sample := range(audio) {
			if t < self.extraSamples {
				for c,amplitude := range(sample) {
					destBuffer[t1][c] = amplitude * d
				}
			} else {
				for c,amplitude := range(sample) {
					destBuffer[t1][c] = amplitude * d + 
				}
			}
		}
	}
	
	// while incoming audio available
		// read through input frame
		// accumulate global sample count
		// if global sample < delayTime in samples
			// copy dry source
		// else
			// mix source w/ delay
		// if global sample count % frameSize == 0
			// send current buffer to to Sink
			// switch current buffer
	// for extra samples
		// write delay data
		// if global sample count % frameSize == 0
			// send current buffer to to Sink
			// switch current buffer
	// for frameSize % extra samples
		// pad with zeros
	
	
}

func (self *DelayFilter) initBuffers() {
	self.bufferSize = self.header.FrameSize
	
	/* initialize floatSamples buffer */
	self.buffers = make([][][]float32, 2)
	for i,_ := range(self.buffers) {
		self.buffers[i] = make([][]float32, self.header.FrameSize)

		for j,_ := range(self.buffers[i]) {
			self.buffers[i][j] = make([]float32, self.header.Channels)
		}
	}
	
	self.frameCopy = make([][]float32, self.header.FrameSize)
	for i,_ := range self.frameCopy {
		self.frameCopy[i] = make([]float32, self.header.Channels)
	}
}
