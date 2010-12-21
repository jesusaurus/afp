// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package delay

import (
	"afp"
	"afp/flags"
	"os"
	"math"
)

type DelayFilter struct {
	context *afp.Context
	header afp.StreamHeader
	samplesPerSecond int
	samplesPerMillisecond int
	delayTimeInMs int
	extraSamples int
	mixBufferSize int64
	bufferSize int32
	channels int16
	bytesPerSample int16
	buffers [][][]float32
	mixBuffer [][]float32
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
/*	var w *int = parser.Int("w", 40, "The wet (delayed) signal ratio: 0 (dry) to 100 (wet)")*/
	parser.Parse()
	
	self.delayTimeInMs = *t	
	
	if self.delayTimeInMs <= 0 {
		return os.NewError("Delay time must be greater than zero")
	}

	return nil
}

func (self *DelayFilter) Start() {
	self.header = <-self.context.HeaderSource
	self.context.HeaderSink <- self.header

	self.samplesPerMillisecond = int(self.header.SampleRate / 1000)
	self.extraSamples = self.delayTimeInMs * self.samplesPerMillisecond;
	
	// the mixBuffer is a ring buffer, each subsection size is self.header.FrameSize, and has n+1 buffers
	// ie, if the delay size is <= frameSize, we have a ring buffer of size 2
	self.mixBufferSize = int64(math.Ceil(float64(self.extraSamples) / float64(self.header.FrameSize))) + 1
	
	self.initBuffers()
	self.process()
}

func (self *DelayFilter) process() {
	var (
		t int64 = 0
		d float32 = 0.75
		w float32 = 0.25
		mbStart int64 = 0
		mbOffset int64 = 0
	)

	d = w
	w = d
	
	for audio := range(self.context.Source) {
		// create a destination buffer
		destBuffer := makeBuffer(self.header.FrameSize, self.header.Channels)
		
		// set mixBuffer to current buffer in the ring to be filled & copy the source audio into that buffer
		mixBuffer := self.mixBuffer[mbStart * int64(self.header.FrameSize):((mbStart+1)*int64(self.header.FrameSize))-1]
		copy(mixBuffer, audio[:])

		println("t: ", t, " mbStart: ", mbStart, " mbOffset: ", mbOffset, " from: ", mbStart * int64(self.header.FrameSize), " to: ", ((mbStart+1)*int64(self.header.FrameSize))-1)

		for t1,sample := range(audio) {
/*			for c,_ := range(sample) {
				(*destBuffer)[t1][c] = self.mixBuffer[mbOffset][c]
			}
			if t > int64(self.extraSamples) {
				mbOffset++
				mbOffset %= (self.mixBufferSize * int64(self.header.FrameSize))
			}
			t++
*/
/*			(*destBuffer)[t1] = sample */

			if t < int64(self.extraSamples) {
				for c,_ := range(sample) {
					(*destBuffer)[t1][c] = 0 * w * d // amplitude * d
				}
			} else {
				for c,_ := range(sample) {
					(*destBuffer)[t1][c] = self.mixBuffer[mbOffset][c] // + amplitude * d
				}
				mbOffset++
			}
			
			if (t == int64(self.extraSamples)) {
				println("Starting delay at ", t, " mbOffset: ", mbOffset)
			}

			if (mbOffset >= (self.mixBufferSize * int64(self.header.FrameSize))) {
				mbOffset = 0
			}
			t++
		
		}
		
/*		self.context.Sink <- mixBuffer */
		self.context.Sink <- *destBuffer
/*		self.context.Sink <- self.mixBuffer[mbStart * int64(self.header.FrameSize):((mbStart+1)*int64(self.header.FrameSize))]*/

		mbStart++
		mbStart %= self.mixBufferSize
		
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

func makeBuffer(size int32, channels int8) *[][]float32 {
	b := make([][]float32, size)
	for i,_ := range(b) {
		b[i] = make([]float32, channels)
	}
	
	return &b
}

func (self *DelayFilter) initBuffers() {
	self.mixBuffer = make([][]float32, self.mixBufferSize * int64(self.header.FrameSize))
	for i,_ := range self.mixBuffer {
		self.mixBuffer[i] = make([]float32, self.header.Channels)
	}
}

func (self *DelayFilter) Stop() os.Error {
	return nil
}