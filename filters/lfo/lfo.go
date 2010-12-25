// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package lfo

import (
	"afp"
	"afp/flags"
	"os"
	"strings"
)

type LFOFilter struct {
	ctx *afp.Context
	amp, freq float32
	mixchan int
	oscGen func(*afp.StreamHeader) func() float32 //Ooogly
}

const ALL = -1

var chStrToVal = map[string]int { "all" : ALL, "left" : 0, "right" : 1 }

var oscStrToVal = map[string]func(*afp.StreamHeader) func() float32 {
	"tri" : getTriangleOscillator,
	"triangle" : getTriangleOscillator,
	"sin" : getSinOscillator,
	"sine" : getSinOscillator,
	"square" : getSquareOscillator,
	"sqr" : getSquareOscillator,
	"sawtoothe" : getSawOscillator,
	"saw" : getSawOscillator,	
}

func (self *LFOFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	parser := flags.FlagParser(args)
	parser.Float32Var(&self.freq, "f", 10, "The frequency of the signal to mix in (in Hz).")
	parser.Float32Var(&self.amp, "a", 0.5, "The amplitude of the signal to mix in.  Between 0 and 1.")
	ch := parser.String("c", "left", "The channel to mix the signal into. May be all, left," +
		" right, or an integer between 0 and the number of channels in the input signal.") 
	shape := parser.String("s", "triangle", "The type of LFO signal to generate: square, sine, sawtoothe, or triangle")
	modTarg := parser.String("m", "volume", "What part of the signal should the LFO modulate: volume, pitch, cutoff")
	parser.Parse()

	if amp < 0 || amp > 1 {
		return os.NewError("Amplitude must be between 0 and 1.")
	}

	chL := strings.ToLower(ch)

	if mc, ok := chArgToVal[chL]; ok {
		self.mixchan = mc
	} else if mc, err := strconv.Atoi(chL); err == nil && mc >= 0 {
		self.mixchan = mc
	} else {
		return os.NewError(fmt.Sprintf("Could not parse '%s' as a valid channel choice." +
			" Must be one of: all, left, right, or an integer indicating the channel to be mixed with.", ch))
	}

	return nil
}

func (self *LFOFilter) Stop() os.Error {
	return nil
}

func (self *LFOFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *LFOFilter) Start() {
	header := <-self.ctx.HeaderSource
	self.ctx.HeaderSink <- header

	osc := self.getTriangleOscillator(&header)

	if self.mixchan == ALL {
		self.mixAll(osc)
	} else {
		self.mixOne(osc)
	}
}

func (self *LFOFilter) mixAll(osc func() float32) {

	for frame := range self.ctx.Source {

		self.ctx.Sink <- frame
	}	
}

func (self *LFOFilter) mixOne(osc func() float32) {

	for frame := range self.ctx.Source {

		self.ctx.Sink <- frame
	}	
}
/**
 * Get 
 */
func (self *LFOFilter) getTriangleOscillator(header *afp.StreamHeader) (func() float32) {
	period := header.SampleRate / self.freq //Roughly, period in slices
	delta := 4 * self.amp / period 
	amp := self.amp //Wary of references to self in a lambda
	var val float32 = 0

	return func() float32 {
		ret := val

		val += delta

		if val >= amp || val <= -amp {
			delta = -delta
		} 

		return ret
	}

}

func NewLFO() afp.Filter {
	return &LFOFilter{}
}

