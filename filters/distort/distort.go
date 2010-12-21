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
    gain, clip, hardness Float32
    clipper func(*DistortFilter)
}

var clipTypes = map[string]func(*DistortFilter) {
    "hard" : hard,
	"variable" : variable,
	"cubic" : cubic,
	"foldback" : foldback,
}

func (self *DistortFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

    fParse := flags.FlagParser(args)
	fParse.Float32Var(&self.gain, "g", 1.0,
		"Signal gain to apply before clipping. Must be > 0.")
    fParse.Float32Var(&self.clip, "c", 1.0,
        "The amplitude at which to clip the signal. Must be in (0,1)")
	fParse.Float32Var(&self.hardness, "k", 100,
			"Clipping 'hardness' for the variable clipping filter. Must be" +
			" in [0,\u221E), where 0 is no clipping and \u221E is hard clipping.")
	clipType := fParse.String("t", 
		"soft", "The type of clipping used: hard, variable, cubic, or foldback.")
	
	fParse.Parse()

	if self.gain <= 0 {
		return os.NewError("Gain must be greater than 0.")
	}

	if self.clip > 1 || self.clip < 0 {
		return os.NewError("Clipping level must be between 0 and 1")
	}

	self.clipper, ok := clipTypes[clipType]

	if !ok {
		return os.NewError("Clipping type must be one of: hard, soft, overflow, or foldback")
	}

	if clipper != variable && self.hardness < 0 {
		return os.NewError("Hardness must be in [0,\u221E).")
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
	
	//A single fold may cause the signal to exceed the clip level 
	//on the other side, so we may need to fold multiple times
	for sample > clip || sample < -clip {
		if sample > clip {
			sample = 2 * clip - sample
		} else {
			sample = clip + sample
		}
	}

	return sample
}

func cubic(f *DistortFilter) {
	for frame := range f.ctx.Source {
		for slice := range frame {
			for ch, sample := range slice {
				frame[slice][ch] = cubicClip(sample * f.gain, f.clip)
			}
		}
		self.ctx.Sink <- frame
	}
}


//This algorithm found at:
//https://ccrma.stanford.edu/~jos/pasp/Soft_Clipping.html#29299
//       { -2/3 * clip        x <= -clip
// out = {  x - x^3/3   -clip < x < clip
//       {  2/3 * clip        x >= clip
func cubicClip(sample, clip float32) float32 {
	if sample >= clip {
		sample = 0.66666666666666666666666 * clip
	} else if sample <= -clip {
		sample = -0.66666666666666666666666 * clip
	} else {
		sample = sample - sample * sample * sample / 3
	}
	return sample
}

func variable(f *DistortFilter) {
	for frame := range f.ctx.Source {
		for slice := range frame {
			for ch, sample := range slice {
				frame[slice][ch] = variableClip(f.hardness, sample * f.gain, f.clip)
			}
		}
		self.ctx.Sink <- frame
	}
}

func variableClip(k, gain, sample, clip float32) float32 {
}

//Provides a good enough approximation of atan
//in [-1,1].  
func fastAtan(x float32) float32 {
	return x / (1 + .28 * x * x)
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
