// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

//inspired by http://www.musicdsp.org/

package echo

import (
    "afp"
    "os"
)

type EchoFilter struct {
    //standard filter stuff
    context *afp.Context
    header afp.StreamHeader

    //decay attenuation: between 0 and 1
    decay float32

    //input and output buffers
    drySignal [][]float32
    wetSignal [][]float32
}

func (self *EchoFilter) GetType() int {
    return afp.PIPE_LINK
}

func NewEchoFilter() afp.Filter {
    return &EchoFilter{}
}

func (self *EchoFilter) Usage() {
    //TODO: add usage
}

func (self *EchoFilter) Init(ctx *afp.Context, args []string) os.Error {
    self.context = ctx
    //TODO: add argument parsing for decay rate
    self.decay = .2

    return nil
}

func (self *EchoFilter) Start() {
    self.header = <-self.context.HeaderSource
    self.context.HeaderSink <- self.header

    //delay offsets for 3 reflections
    offset1 := int32(20) //magic number
    offset2 := int32(35) //magic number
    offset3 := int32(42) //magic number

    //make the input buffer twice the frame size
    self.drySignal = <-self.context.Source
    self.drySignal = append(self.drySignal, <-self.context.Source...)
    length := 2 * self.header.FrameSize

    //a couple of empty buffers with the same dimensions as our input signal
    var zero []float32 = make([]float32, self.header.Channels)
    var zeros [][]float32
    for i := int32(0); i < self.header.FrameSize; i++ {
        zeros = append(zeros, zero)
    }

    self.wetSignal = makeBuffer(self.header.FrameSize*2, self.header.Channels)

    for nextFrame := range self.context.Source {

        outBuffer := makeBuffer(self.header.FrameSize, self.header.Channels)

        for i := int32(0); i < self.header.FrameSize; i++ {
            for j := int8(0); j < self.header.Channels; j++ {
                self.wetSignal[i][j] = self.drySignal[i][j]
                self.wetSignal[i+offset1][j] = self.wetSignal[i+offset1][j] + (self.drySignal[i][j] * self.decay )
                self.wetSignal[i+offset2][j] = self.wetSignal[i+offset2][j] + (self.drySignal[i][j] * self.decay )
                self.wetSignal[i+offset3][j] = self.wetSignal[i+offset3][j] + (self.drySignal[i][j] * self.decay )
            }
        }

        outBuffer = self.wetSignal[0:self.header.FrameSize]
        self.context.Sink <- outBuffer

        self.wetSignal = self.wetSignal[self.header.FrameSize:]
        self.wetSignal = append(self.wetSignal, zeros...)

        self.drySignal = self.drySignal[self.header.FrameSize:]
        self.drySignal = append(self.drySignal, nextFrame...)
    }

    //TODO: pad with silence

    //flush the signals
    for i := int32(0); i < length; i++ {

        outBuffer := makeBuffer(self.header.FrameSize, self.header.Channels)

        //apply echo/reverb
        for j := int8(0); j < self.header.Channels; j++ {
            self.wetSignal[i][j] = (self.wetSignal[i][j]  ) + (self.drySignal[i][j] * self.decay )
            self.wetSignal[i+offset1][j] = (self.wetSignal[i][j]  ) + (self.drySignal[i][j] * self.decay )
            self.wetSignal[i+offset2][j] = (self.wetSignal[i][j]  ) + (self.drySignal[i][j] * self.decay )
            self.wetSignal[i+offset3][j] = (self.wetSignal[i][j]  ) + (self.drySignal[i][j] * self.decay )
        }

        //wrap
        if i == self.header.FrameSize {
            outBuffer = self.wetSignal[0:self.header.FrameSize]
            self.context.Sink <- outBuffer

            outBuffer := makeBuffer(self.header.FrameSize, self.header.Channels)

            self.wetSignal = self.wetSignal[self.header.FrameSize:]
            self.drySignal = self.drySignal[self.header.FrameSize:]

            i = 0
            length -= self.header.FrameSize
        }
    }
}

func (self *EchoFilter) Stop() os.Error {
    //TODO
    return nil
}

// allocate a buffer for samples
func makeBuffer(size int32, channels int8) [][]float32 {
	b := make([][]float32, size)
	for i, _ := range b {
		b[i] = make([]float32, channels)
	}

	return b
}

