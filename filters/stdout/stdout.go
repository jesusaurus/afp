// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package stdout;

import (
	"afp"
	"os"
	"encoding/binary"
)

type StdoutSink struct {
	context *afp.Context
	header afp.StreamHeader
}

func NewStdoutSink() afp.Filter {
	return &StdoutSink{}
}

func (self *StdoutSink) Init(ctx *afp.Context, args []string) os.Error {
    self.context = ctx
	return nil
}

func (self *StdoutSink) GetType() int {
    return afp.PIPE_SINK
}

func (self *StdoutSink) Start() {
    self.header = <-self.context.HeaderSource
    ibuf := make([]int16, int32(self.header.Channels) * self.header.FrameSize)

    for buffer := range self.context.Source { //reading a [][]float32
        length := int(self.header.FrameSize)
        chans := int(self.header.Channels)

		streamOffset := 0
        //interleave the channels
        for i := 0; i < length; i ++ {
            for j := 0; j < chans; j++ {
                ibuf[streamOffset] = int16(buffer[i][j] * 32767.0)
				streamOffset += 1
            }
        }

		err := binary.Write(os.Stdout, binary.LittleEndian, ibuf)
		if (err != nil) {
			os.Stderr.WriteString(err.String())
		}
    }

    return
}

func (self *StdoutSink) Stop() os.Error {
    return nil
}

