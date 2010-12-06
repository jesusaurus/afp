// Copyright (c) 2010 Go Fightclub Authors

////
// Chiptune-ify
// This filter will make your music sound like it came from an 8-bit video game

package chip

import (
    "afp"
)

type ChiptuneFilter struct {
    ctx *afp.Context
    header StreamHeader
}

func (self *ChiptuneFilter) GetType() int {
    return afp.PIPE_LINK
}

func (self *ChiptuneFilter) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx
    self.header <-self.ctx.HeaderSource
    self.ctx.HeaderSink <- self.header
    return nil
}

func (self *ChiptuneFilter) Start() {
    for buffer := range self.ctx.Source {
        var samples [][]float32 //reversed dimensions
        samples = make([][]float32, self.header.Channels)
        length = len(buffer)

        for channel := 0; channel < self.ctx.Channels; channel++ {
            samples[channel] = make([]float32, length)
        }

        for channel, i := 0, 0; i < length; channel, i = channel + 1, i + self.header.Channels {
        }
    }
    return
}

func (self *ChiptuneFilter) Stop() os.Error {
    return nil
}
