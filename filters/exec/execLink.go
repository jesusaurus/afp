// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package fexec

import (
	"afp"
)

type ExecLink struct {
	execFilter
}

func NewExecLink() afp.Filter {
	return &ExecLink{execFilter{}}
}

func (self *ExecLink) GetType() int {
	return afp.PIPE_LINK
}

func (self *ExecLink) Start() {

	go self.encoder()
	go self.decoder()

	if self.context.Verbose {
		go self.errors()
	}

	self.wait()
}
