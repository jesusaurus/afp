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
    context *afp.Context
    header afp.StreamHeader
    decay float32 //decay factor: between 0 and 1
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
    offset1 := 20 //magic number
    offset2 := 35 //magic number
    offset3 := 42 //magic number

    drySignal := <-self.context.Source //[][]float32
    frameSize := len(drySignal)
    //make the dry signal buffer twice the frame size
    drySignal = append(drySignal, <-self.context.Source...)
    length := len(drySignal)

    var (
        wetSignal [][]float32
    )

    //a couple of empty buffers
    var zero []float32
    for _, _ = range drySignal[0] {//we don't care about the data, just the dimensions
        zero = append(zero, 0)
    }
    var zeros [][]float32
    for i := 0; i < frameSize; i++ {
        zeros = append(zeros, zero)
    }

    for i := 0; i < 3; i++ {
        //make our buffers 3 frames large
        wetSignal = append(wetSignal, zeros...)
    }

    for nextFrame := range self.context.Source {

        for i := 0; i < frameSize; i++ {
            for j := int8(0); j < self.header.Channels; j++ {
                //ECHO, Echo, echo...
                wetSignal[i][j] = (wetSignal[i][j] ) + (drySignal[i][j] * self.decay )
                wetSignal[i+offset1][j] = (wetSignal[i][j]  ) + (drySignal[i][j] * self.decay )
                wetSignal[i+offset2][j] = (wetSignal[i][j]  ) + (drySignal[i][j] * self.decay )
                wetSignal[i+offset3][j] = (wetSignal[i][j]  ) + (drySignal[i][j] * self.decay )
            }
        }

        self.context.Sink <- wetSignal[0:frameSize]
        wetSignal = wetSignal[frameSize:]
        wetSignal = append(wetSignal, zeros...)

        drySignal = drySignal[frameSize:]
        drySignal = append(drySignal, nextFrame...)
    }

    //TODO: pad with silence

    //flush the signals
    for i := 0; i < length; i++ {
        //apply echo/reverb
        for j := int8(0); j < self.header.Channels; j++ {
            wetSignal[i][j] = (wetSignal[i][j]  ) + (drySignal[i][j] * self.decay )
            wetSignal[i+offset1][j] = (wetSignal[i][j]  ) + (drySignal[i][j] * self.decay )
            wetSignal[i+offset2][j] = (wetSignal[i][j]  ) + (drySignal[i][j] * self.decay )
            wetSignal[i+offset3][j] = (wetSignal[i][j]  ) + (drySignal[i][j] * self.decay )
        }

        //wrap
        if i == frameSize {
            self.context.Sink <- drySignal[0:frameSize]
            wetSignal = wetSignal[frameSize:]
            drySignal = drySignal[frameSize:]
            i = 0
            length -= frameSize
        }
    }
}

func (self *EchoFilter) Stop() os.Error {
    //TODO
    return nil
}
