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
