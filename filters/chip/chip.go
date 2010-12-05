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

    header <-ctx.HeaderSource

    return nil
}

func (self *ChiptuneFilter) Start() {
    return
}
