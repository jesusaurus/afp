package fexec

import (
	"afp"
	"os"
)

type ExecSink struct {
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

	if self.Verbose {
		go self.errors()
	}
}

func Stop() os.Error {

	return nil
}