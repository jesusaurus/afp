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

	if *linear == *exp {
		return os.NewError("You must specify exactly one convolution algorithm. " +
			" Available choices are: linear, exp")
	}

	if *linear {
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
//If O_k is the kth sample in the output and I_k the kth symbol in the input,
//then each sample in the input will equal:
// O_n = I_{2n - 1}/4 + I_{2n}/2 + I_{2n + 1}
//where I_-1 refers to the last sample in the previous frame  
func linearDS(carryOver []float32, input [][]float32) [][]float32 {
    var outSample float32
	output := input[:len(input) / 2]

    for in, out := 0, 0; in < len(input); out++ {
		for j := range input[in] {
			outSample = carryOver[j] + input[in][j] / 2
			carryOver[j] = input[in + 1][j] / 4;
			output[out][j] = outSample + carryOver[j]
		}
		in += 2
	}
	return output
}

//This algorithm adapted from mumart[AT]gmail[DOT]com
//Found at http://www.musicdsp.org/showArchiveComment.php?ArchiveID=214
//We mix an exponentially decreasing fraction of each sample into every one following
//So, the O_n, the nth sample in the output, may be expressed as:
// O_n = sum I_k / 2^(2n - k + 1) for k = 0 to 2n
//Where I_k is the kth sample in the input
func expDS(carry []float32, input [][]float32) [][]float32 {
	output := input[:len(input) / 2]

    for in, out := 0, 0; in < len(input); out++ {
		for ch := range input[in] {
			carry[ch] = (carry[ch] + input[in][ch]) / 2;
			output[out][ch] = carry[ch] 
			carry[ch] = (carry[ch] + input[in + 1][ch] ) / 2;
		}
		in += 2;
    }
	return output
}
