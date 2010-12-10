// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package portaudio

import (
	"afp"
	"os"
	"fmt"
)

// #include "pasink.c"
import "C"

type PASink struct {
	context *afp.Context
	header afp.StreamHeader
	output_data C.pa_output_data
}

func NewPASink() afp.Filter {
    return &PASink{}
}

func (self *PASink) GetType() int {
    return afp.PIPE_SINK
}

func (self *PASink) Init(ctx *afp.Context, args []string) os.Error {
    self.context = ctx
	
    return nil
}

func (self *PASink) Start() {
    self.header = <-self.context.HeaderSource

	err := C.init_portaudio_output(C.int(self.header.Channels), C.int(self.header.SampleRate), C.int(self.header.FrameSize), &self.output_data)
    if (err != 0) {
		os.Stderr.WriteString("Problem!")
        panic(os.NewError(fmt.Sprintf("Initialize portaudio failed, error: %d", err)))
    }

    cbuf := make([]float32, int32(self.header.Channels) * self.header.FrameSize)

    for buffer := range self.context.Source { //reading a [][]float32
        length := int(self.header.FrameSize)
        chans := int(self.header.Channels)

		streamOffset := 0
        //interleave the channels
        for i := 0; i < length; i ++ {
            for j := 0; j < chans; j++ {
                cbuf[streamOffset] = buffer[i][j]
				streamOffset++
            }
        }

        //write some data to portaudio
        C.send_output_data((*C.float)(&cbuf[0]), &self.output_data, 0)
    }

	// terminate the stream 
	C.send_output_data((*C.float)(&cbuf[0]), &self.output_data, 1)

    return
}

func (self *PASink) Stop() os.Error {
	os.Stderr.WriteString("Stopping.")
	C.close_portaudio(&self.output_data)
    return nil
}
