// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package lfo

import (
	"afp"
	"afp/flags"
	"os"
)

type LFOFilter struct {
	ctx *afp.Context
}

func (self *LFOFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	parser := flags.FlagParser(args)
	f := parser.Float32("f", 10, "The frequency of the signal to mix in (in Hz).")
	a := parser.Float32("a", 0.5, "The amplitude of the signal to mix in.  Between 0 and 1.")
	parser.Parse()

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

	for frame := range self.ctx.Source {

		self.ctx.Sink <- frame
	}
}

func NewLFO() afp.Filter {
	return &LFOFilter{}
}

