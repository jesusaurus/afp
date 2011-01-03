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

    drySignal := <-self.context.Source
    frameSize := len(drySignal)
    //make the dry signal buffer twice the frame size
    drySignal = append(drySignal, <-self.context.Source...)
    length := len(drySignal)

    reflect1 := make([][]float32, length)
    reflect2 := make([][]float32, length)
    reflect3 := make([][]float32, length)
    wetSignal := make([][]float32, length)

    //a couple of empty buffers
    var zero float32[]
    for _, _ := range drySignal[0] {//we don't care about the data, just the dimensions
        zero = append(zero, 0)
    }
    var zeros float32[][]
    for i := 0; i < frameSize; i++ {
        zeros = append(zeros, zero)
    }

    //initialize our reflection buffers with silence
    for i := 0; i < offset1; i++ {
        reflect1[i] = zero
        reflect2[i] = zero
        reflect3[i] = zero
    } for i < offset2 {
        reflect2[i] = zero
        reflect3[i] = zero
        i++
    } for i < offset3 {
        reflect3[i] = zero
        i++
    }

    for nextFrame := range self.context.Source {

        for i := 0; i < frameSize; i++ {
            for j := int8(0); j < self.header.Channels; j++ {
                //ECHO, Echo, echo...

                reflect1[i+offset1][j] = drySignal[i][j] * self.decay
                reflect2[i+offset2][j] = drySignal[i][j] * self.decay
                reflect3[i+offset3][j] = drySignal[i][j] * self.decay

                wetSignal[i][j] = reflect1[i] + reflect2[i] + reflect3[i]
            }
        }

        self.context.Sink <- wetSignal[0:frameSize]
        wetSignal = wetSignal[frameSize:]
        wetSignal = append(wetSignal, zeros)

        drySignal = drySignal[frameSize:]
        drySignal = append(drySignal, nextFrame...)
    }

    //TODO: pad with silence

    //flush the signals
    for i := 0; i < length; i++ {
        //apply echo/reverb
        for j := int8(0); j < self.header.Channels; j++ {
            reflect1[i+offset1][j] = drySignal[i][j] * self.decay
            reflect2[i+offset2][j] = drySignal[i][j] * self.decay
            reflect3[i+offset3][j] = drySignal[i][j] * self.decay

            wetSignal[i][j] = reflect1[i] + reflect2[i] + reflect3[i]
        }

        //wrap
        if i == frameSize {
            self.context.Sink <- wetSignal[0:frameSize]
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
