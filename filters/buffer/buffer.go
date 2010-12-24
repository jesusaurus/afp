// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package buffer

import (
	"afp"
	"afp/lib/util"
	"afp/flags"
	"os"
)

type BufferFilter struct {
	ctx *afp.Context
	toBuffer int
}

func (self *BufferFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	parser := flags.FlagParser(args)
	parser.IntVar(&self.toBuffer, "b", afp.CHAN_BUF_LEN, "The number of frames to buffer.")
	parser.Parse()

	return nil
}

func (self *BufferFilter) Stop() os.Error {
	return nil
}

func (self *BufferFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *BufferFilter) Start() {
	header := <-self.ctx.HeaderSource
	self.ctx.HeaderSink <- header

	source := util.Buffer(self.toBuffer, self.ctx.Source)

	for frame := range source {
		self.ctx.Sink <- frame
	}
}

func NewBuffer() afp.Filter {
	return &BufferFilter{}
}

