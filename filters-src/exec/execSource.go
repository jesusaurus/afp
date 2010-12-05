package fexec

import (
	"afp"
	//"os"
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

