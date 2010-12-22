// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package fexec

import (
	"afp"
)

type ExecSource struct {
	execFilter
}

func NewExecSource() afp.Filter {
	return &ExecSource{execFilter{}}
}

func (self *ExecSource) GetType() int {
	return afp.PIPE_SOURCE
}

func (self *ExecSource) Start() {

	go self.decoder()

	if self.context.Verbose {
		go self.errors()
	}

	self.wait()
}
