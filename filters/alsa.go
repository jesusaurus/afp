// Copyright (c) 2010 Go Fightclub Authors

package alsa

import (
  "afp"
)

type AlsaFilter struct {
    ctx *afp.Context
    header StreamHeader
}

func (self *AlsaFilter) GetType() int {
    return afp.PIPE_SINK
}

func (self *AlsaFilter) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx

    header <-ctx.HeaderSource
}

func (self *AlsaFilter) Start() {

}
