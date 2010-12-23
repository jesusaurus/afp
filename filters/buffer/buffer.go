// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

//This is not a legal Go program, rather it provides a skeletal
//filter to serve as a minimal base for developing filters.

package buffer

import (
	"afp"
	"afp/lib/util"
	"os"
)

type BufferFilter struct {
	ctx *afp.Context
}

func (self *BufferFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	return nil
}

func (self *BufferFilter) Stop() os.Error {
	return nil
}

func (self *BufferFilter) GetType() int {
	return afp.PIPE_< SOURCE | LINK | SINK >
}

func (self *BufferFilter) Start() {
	header := <-self.ctx.HeaderSource
	self.ctx.HeaderSink <- header

	for frame := range self.ctx.Source {

	}
}

func NewBuffer() afp.Filter {
	return &BufferFilter{}
}

