// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package echo

import (
    "afp"
    "flags"
    "os"
)

type EchoFilter struct {
    context *afp.Context
    header afp.StreamHeader
}

func (self *EchoFilter) GetType() int {
    return afp.PIPE_LINK
}

func NewEchoFilter() afp.Filter {
    return &EchoFilter{}
}

func (self *EchoFilter) Usage() {
    
}

func (self *EchoFilter) Init(ctx *afp.Context, args []string) os.Error {
    self.context = ctx
    return nil
}

func (self *EchoFilter) Start() {
    
}

func (self *EchoFilter) Stop() os.Error {
    //TODO
    return nil
}
