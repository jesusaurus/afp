// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package distort

import (
	"afp"
	"afp/flags"
	"math"
)

type DistortFilter struct {
	ctx *afp.Context
    gain, clip Float32
    clipper func(*DistortFilter)
}

var clipTypes = map[string]func(*DistortFilter) {
    "hard" : hard,
	"soft" : nil,
	"overflow" : nil,
	"foldback" : foldback,

}

func (self *DistortFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

    fParse := flags.FlagParser(args)
	fParse.Float32Var("g", 1.0,
		"Signal gain to apply before clipping. Must be greater than 0.")
    clipLevel := fParse.Float64("c", 1.0,
        "The amplitude at which to clip the signal. Must be between 0 and 1.")
    clipType := fParse.String("t", 
		"soft", "The type of clipping used: hard, soft, overflow, or foldback.")
	
	fParse.Parse()

	if gain64 <= 0 {
		return os.NewError("Gain must be greater than 0.")
	}
	self.gain = float32(gain64)

	if clip64 > 1 || clip64 < 0{
		return os.NewError("Clipping level must be between 0 and 1")
	}
	self.clip = float32(clip64)
	self.clipper, ok := clipTypes[clipType]

	if !ok {
		return os.NewError("Clipping type must be one of: hard, soft, overflow, or foldback")
	}

	return nil
}

func (self *DistortFilter) Stop() os.Error {
	return nil
}

func (self *DistortFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *DistortFilter) Start() {
	self.ctx.HeaderSink <- (<-self.ctx.HeaderSource)
	self.clipper(self)
}

func hard(f *DistortFilter) {
	for frame := range f.ctx.Source {
		for slice := range frame {
			for ch, sample := range slice {
				frame[slice][ch] = hardMin(f.clip, sample * f.gain)
			}
		}
		self.ctx.Sink <- frame
	}
}

//Min function which knows about hard(). 
//specifically that clip will always be positive
func hardMin(clip, sprime float32) {
	var t float32

	if sprime < 0 {
		t = -sprime
	} else {
		t = sprime
	}
	
	if t > clip {
		return clip
	}

	return sprime
}

func foldback(f *DistortFilter) {
	for frame := range f.ctx.Source {
		for slice := range frame {
			for ch, sample := range slice {
				frame[slice][ch] = fold(sample * f.gain, f.clip)
			}
		}
		self.ctx.Sink <- frame
	}
}

//Helper function for foldback
//Computes the actual value of a sample
func fold(sample, clip float32) float32 {	
	for sample > clip || sample < -clip {
		if sample > clip {
			sample = 2 * clip - sample
		} else {
			sample = clip + sample
		}
	}
	return sample
}

//Original C version by Alexander Kritov 
//http://www.musicdsp.org/archive.php?classid=1#68  
func soft(x, a, N, fi float32) {
    var (
        s1 = pow(a, N-1.0) * sin((N - 1.0) * x + fi)
        s2 = pow(a, N) * sin(N * x + fi)
        s3 = a * sin(x + fi)
        s4 = 1.0 - (2 * a * cos(x)) + (a * a)
    )

    if s4 == 0 {
        return 0;
    } else {
        return (sin(fi) - s3 - s2 +s1) / s4;
    }
}
