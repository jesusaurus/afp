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
	downsampler func([]float32, [][]float32) [][]float32
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
	
	return nil
}

func (self *Halfsampler) Stop() os.Error {
	return nil
}

func (self *Halfsampler) GetType() int {
	return afp.PIPE_LINK
}

func (self *Halfsampler) Start() {
	header := <-self.ctx.HeaderSource

	//Do we need to be slightly more clever here?
	headerCopy := header
	headerCopy.SampleRate = header.SampleRate / 2
	headerCopy.FrameSize = header.FrameSize / 2

	//We may not be able to know resulting content length
	//Is this necessarily true?
	headerCopy.ContentLength = 0 

	self.ctx.HeaderSink <- headerCopy

	carryOver := make([]float32, header.Channels)

	for frame := range self.ctx.Source {
		self.ctx.Sink <- self.downsampler(carryOver, frame)
	}
}

func NewHalfsampler() afp.Filter {
	return &Halfsampler{}
}

//This algorithm adapted from mumart[AT]gmail[DOT]com
//Found at http://www.musicdsp.org/showArchiveComment.php?ArchiveID=214
func linearDS(carryOver []float32, input [][]float32) [][]float32 {
    var outSample float32
	outBuff := input[:len(input / 2)]

    for i, outInd := 0, 0; i < len(input); outInd++ {
		for j := range input[i] {
			output_sam =  carryOver[j] + input[i][j] / 2
			i++
			
			carryOver[j] = input[i][j] / 2;
			i++
			
			outBuff[outInd][j] = output_sam + carryOver[j]
		}
	}
}

//This algorithm adapted from mumart[AT]gmail[DOT]com
//Found at http://www.musicdsp.org/showArchiveComment.php?ArchiveID=214
func expDS( int *input_buf, int *output_buf, int output_count ) {
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