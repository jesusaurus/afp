// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package halfsample

import (
	"afp"
	"afp/flags"
	"os"
)

type Halfsampler struct {
	ctx *afp.Context
	downsampler func([][]float32) [][]float32
}

func (self *Halfsampler) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	parser := flags.FlagParser(args)
	linear := parser.Bool("linear", false, "Use a linear convolution before downsampling")
	exp := parser.Bool("exp", false, "Use an exponential convolution before downsampling")
	parser.Parse()

	if linear == exp {
		return os.NewError("You must specify exactly one convolution algorithm. " +
			" Available choices are: linear, exp")
	}

	if linear {
		self.downsampler = linearDS
	} else {
		self.downsampler = expDS
	}
}

func (self *Halfsampler) Stop() os.Error {
	return nil
}

func (self *Halfsampler) GetType() int {
	return afp.PIPE_LINK
}

func (self *Halfsampler) Start() {
	header := <-self.ctx.HeaderSource
	headerCopy := header
	headerCopy.SampleRate = header.SampleRate / 2
	headerCopy.FrameSize = header.FrameSize / 2
	self.ctx.HeaderSink <- headerCopy

	//Then process the content til there's no more to be had
	for frame := range self.ctx.Source {
		//Process frame
	}
}

func NewHalfsampler() afp.Filter {
	return &Halfsampler{}
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

/*
int filter_state;

void downsample( int *input_buf, int *output_buf, int output_count ) {
    int input_idx, output_idx, input_ep1;
    output_idx = 0;
    input_idx = 0;
    input_ep1 = output_count * 2;
    while( input_idx < input_ep1 ) {
        filter_state = ( filter_state + input_buf[ input_idx ] ) >> 1;
        output_buf[ output_idx ] = filter_state;
        filter_state = ( filter_state + input_buf[ input_idx + 1 ] ) >> 1;
        input_idx += 2;
        output_idx += 1;
    }
}
*/