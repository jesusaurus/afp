// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package distort

import (
	"afp"
	"afp/flags"
	"os"
)

type DistortFilter struct {
	ctx                  *afp.Context
	gain, clip, hardness float32
	clipper              func(*DistortFilter)
}

func NewFilter() afp.Filter {
	return &DistortFilter{}
}

var clipTypes = map[string]func(*DistortFilter){
	"hard":     hard,
	"variable": variable,
	"cubic":    cubic,
	"foldback": foldback,
}

func (self *DistortFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

	fParse := flags.FlagParser(args)
	fParse.Float32Var(&self.gain, "g", 1.0,
		"Signal gain to apply before clipping. Must be > 0.")
	fParse.Float32Var(&self.clip, "c", 1.0,
		"The amplitude at which to clip the signal. Must be in (0,1)")
	fParse.Float32Var(&self.hardness, "h", 10,
		"Clipping 'hardness' for the variable clipping filter. Must be"+
			" in [1,\u221E), where 1 is soft clipping and \u221E is hard clipping.")
	clipType := fParse.String("t",
		"cubic", "The type of clipping used: hard, variable, cubic, or foldback."+
			" See the afp(1) manpage for more info")

	fParse.Parse()

	if self.gain <= 0 {
		return os.NewError("Gain must be greater than 0.")
	}

	if self.clip > 1 || self.clip < 0 {
		return os.NewError("Clipping level must be between 0 and 1")
	}

	tempClipper, ok := clipTypes[*clipType]

	if !ok {
		return os.NewError("Clipping type must be one of: hard, soft, overflow, or foldback")
	}

	self.clipper = tempClipper

	if self.clipper != variable && self.hardness < 1 {
		return os.NewError("Hardness must be in [1,\u221E).")
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
			for ch, sample := range frame[slice] {
				frame[slice][ch] = hardMin(f.clip, sample * f.gain)
			}
		}
		f.ctx.Sink <- frame
	}
}

//Min function which knows about hard(). 
//specifically that clip will always be positive
func hardMin(clip, sprime float32) float32 {
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
			for ch, sample := range frame[slice] {
				frame[slice][ch] = fold(sample*f.gain, f.clip)
			}
		}
		f.ctx.Sink <- frame
	}
}

//Helper function for foldback
//Computes the actual value of a sample
func fold(sample, clip float32) float32 {

	//A single fold may cause the signal to exceed the clip level 
	//on the other side, so we may need to fold multiple times
	for sample > clip || sample < -clip {
		if sample > clip {
			sample = 2*clip - sample
		} else {
			sample = clip + sample
		}
	}

	return sample
}


func cubic(f *DistortFilter) {
	for frame := range f.ctx.Source {
		for slice := range frame {
			for ch, sample := range frame[slice] {
				frame[slice][ch] = cubicClip(sample*f.gain, f.clip)
			}
		}
		f.ctx.Sink <- frame
	}
}

//This algorithm is an adaptation of the one found at:
//https://ccrma.stanford.edu/~jos/pasp/Soft_Clipping.html#29299
//       { -2/3 * clip           x <= -clip
// out = {  x - x^3/3    -clip < x < clip
//       {  2/3 * clip           x >= clip
func cubicClip(sample, clip float32) float32 {
	if sample >= clip {
		sample = 0.66666666666666666666666 * clip
	} else if sample <= -clip {
		sample = -0.66666666666666666666666 * clip
	} else {
		sample = sample - sample*sample*sample/3
	}
	return sample
}

//Variable distortion via a modification of the formula
//from http://www.musicdsp.org/showone.php?id=104
//by scoofy[AT]inf[DOT]elte[DOT]hu
//For each sample, we evaluate:
// c/atan(s) * atan(x*s), where c is the clip level, s is
// the hardness, and x is the sample data.
func variable(f *DistortFilter) {
	//Precompute what we can..
	hardnessMult := f.clip / atan(f.hardness)

	for frame := range f.ctx.Source {
		for slice := range frame {
			for ch, sample := range frame[slice] {
				frame[slice][ch] = hardnessMult * atan(sample*f.hardness)
			}
		}
		f.ctx.Sink <- frame
	}
}

//Provides a good enough approximation of atan
//in [-2,2].  Thanks to antiprosynthesis[AT]hotmail[DOT]com
//Not used at this time.
func fastAtan(x float32) float32 {
	return x / (1 + .28*x*x)
}
