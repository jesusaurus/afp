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
}

var chArgToVal = map[string]int { "all" : -1, "left" : 0, "right" : 1 }

func (self *LFOFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	parser := flags.FlagParser(args)
	parser.Float32Var(&self.freq, "f", 10, "The frequency of the signal to mix in (in Hz).")
	parser.Float32Var(&self.amp, "a", 0.5, "The amplitude of the signal to mix in.  Between 0 and 1.")
	ch := parser.String("c", "left", "The channel to mix the signal into. May be all, left," +
		" right, or an integer between 0 and the number of channels in the input signal.") 
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
	
}

func (self *LFOFilter) mixAll() {
	header := <-self.ctx.HeaderSource
	self.ctx.HeaderSink <- header

	for frame := range self.ctx.Source {

		self.ctx.Sink <- frame
	}
	
}

func (self *LFOFilter) mixOne() {
	header := <-self.ctx.HeaderSource
	self.ctx.HeaderSink <- header

	for frame := range self.ctx.Source {

		self.ctx.Sink <- frame
	}
	
}

func NewLFO() afp.Filter {
	return &LFOFilter{}
}

