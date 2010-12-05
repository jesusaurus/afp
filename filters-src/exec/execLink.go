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
