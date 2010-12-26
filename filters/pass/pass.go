// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package pass

import (
	"afp"
	"afp/flags"
	"afp/fftw"
	"afp/matrix"
	"afp/window"
	"os"
)

const DEBUG = true

type LowPassFilter struct {
	context               *afp.Context
	header                afp.StreamHeader
	cutoffFrequency		  float32
}

func NewLowPassFilter() afp.Filter {
	return &LowPassFilter{}
}

func (self *LowPassFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *LowPassFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.context = ctx

	parser := flags.FlagParser(args)
	var f *float = parser.Float("f", 440, "The cutoff frequency")
	parser.Parse()

	self.cutoffFrequency = float32(*f)

	if self.cutoffFrequency < 0 {
		return os.NewError("The cutoff frequency must be greater than zero")
	}

	return nil
}

func (self *LowPassFilter) Start() {
	self.header = <-self.context.HeaderSource
	self.context.HeaderSink <- self.header

	self.process()
}

func (self *LowPassFilter) process() {
	N := int(self.header.FrameSize * int32(self.header.Channels))
	
	b := make([][]float32, 2)
	b[0] = make([]float32, 0)
	b[1] = make([]float32, 0)
	bn := 0
	
	w0, w1, wd := 0, int(self.header.FrameSize), int(self.header.FrameSize / 2)
	
	// loop over all input data
	for audio := range self.context.Source {

		interleaved := matrix.Interleave(audio)
		b[1 - bn] = append(b[bn], interleaved...)
		bn = 1 - bn
		
		window := window.Hann(b[bn][w0:w1-1])
		
		println("Windowed:")
		for _, amp := range(window) {
			print(amp, " ")
		}
		println(); println();
		
		_ = fftw.RealToReal1D_32(window, true, 3, fftw.MEASURE, fftw.R2HC)

		println("Spectral:")
		for f, _ := range(window) {
			print(window[f], " ")
		}
		println(); println();
		
		_ = fftw.RealToReal1D_32(window, true, 3, fftw.MEASURE, fftw.HC2R)
		
		println("Temporal:")
		for t, amp := range(window) {
			interleaved[t] = amp/float32(N)
			print(interleaved[t], " ")
		}
		println(); println();

		deinterleaved := matrix.Deinterleave(interleaved, int(self.header.FrameSize), int(self.header.Channels))
		
		// increment window interval
		w0, w1 = w0 + wd, w1 + wd

		// send the mixed audio down the pipe
		self.context.Sink <- deinterleaved
	}

}

func (self *LowPassFilter) Stop() os.Error {
	return nil
}

