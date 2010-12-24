// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package pass

import (
	"afp"
	"afp/flags"
	"afp/fftw"
	"afp/matrix"
	"os"
)

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
	// loop over all input data
	for audio := range self.context.Source {

		interleaved := matrix.Interleave(audio)
		
		_ = fftw.RealToReal1D_32(interleaved, true, 3, fftw.MEASURE, fftw.R2HC)
		
/*		for f, _ := range(interleaved) {
			if f < 60 {
				interleaved[f] = 0
			}
		}
*/		
		_ = fftw.RealToReal1D_32(interleaved, true, 3, fftw.MEASURE, fftw.HC2R)
		
		for t, amp := range(interleaved) {
			interleaved[t] = amp/float32(512 * len(audio) * len(audio[0]))
/*			print(interleaved[t], " ")*/
		}

		deinterleaved := matrix.Deinterleave(interleaved, len(audio), len(audio[0]))

		// send the mixed audio down the pipe
		self.context.Sink <- deinterleaved
	}

}

func (self *LowPassFilter) Stop() os.Error {
	return nil
}
