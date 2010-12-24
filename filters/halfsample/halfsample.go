// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package downsample

import (
	"afp"
	"afp/flags"
	"os"
)

type Downsampler struct {
	ctx *afp.Context
}

func (self *Downsampler) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	return nil
}

func (self *Downsampler) Stop() os.Error {
	return nil
}

func (self *Downsampler) GetType() int {
	return afp.PIPE_< SOURCE | LINK | SINK >
}

func (self *Downsampler) Start() {
	//The first thing Start should do is store
	//and pass on the header info.
	header := <-self.ctx.HeaderSource
	self.ctx.HeaderSink <- header

	//Then process the content til there's no more to be had
	for frame := range self.ctx.Source {
		//Process frame
	}
}

func NewSkeleton() afp.Filter {
	return &Downsampler{}
}
/*
int filter_state;

void downsample( int *input_buf, int *output_buf, int output_count ) {
    int input_idx, input_end, output_idx, output_sam;
    input_idx = output_idx = 0;
    input_end = output_count * 2;
    while( input_idx < input_end ) {
        output_sam = filter_state + ( input_buf[ input_idx++ ] >> 1 );
        filter_state = input_buf[ input_idx++ ] >> 2;
        output_buf[ output_idx++ ] = output_sam + filter_state;
    }
}
*/
