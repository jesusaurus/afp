// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package distort

import (
	"afp"
	"math"
)

type DistortFilter struct {
	ctx *afp.Context
}

func (self *DistortFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx
	return nil
}

func (self *DistortFilter) Stop() os.Error {
	return nil
}

func (self *DistortFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *DistortFilter) Start() {
}