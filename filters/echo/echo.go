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

    buffer := <-self.context.Source
    frameSize := len(buffer)
    //make the buffer 2x the frame size
    buffer = append(buffer, <-self.context.Source...)

    for nextFrame := range self.context.Source {
        for i := 0; i < frameSize; i++ {
            for j := int8(0); j < self.header.Channels; j++ {
                //ECHO, Echo, echo...
                buffer[i+1][j] += buffer[i][j] * self.decay
            }
        }

        self.context.Sink <- buffer[0:frameSize]
        buffer = buffer[frameSize:]
        buffer = append(buffer, nextFrame...)
    }
}

func (self *EchoFilter) Stop() os.Error {
    //TODO
    return nil
}
