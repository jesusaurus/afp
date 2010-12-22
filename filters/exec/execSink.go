// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package fexec

import (
	"afp"
)

type ExecSink struct {
	execFilter
}

func NewExecSink() afp.Filter {
	return &ExecSource{execFilter{}}
}

func (self *ExecSink) GetType() int {
	return afp.PIPE_SINK
}

func (self *ExecSink) Start() {

	go self.encoder()

	if self.context.Verbose {
		go self.errors()
	}

	self.wait()
}
