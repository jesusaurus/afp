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
    "hard" : hardCutoff,
	"soft" : nil,
	"overflow" : nil,
	"foldback" : nil,

}

func (self *DistortFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

    fParse := flags.FlagParser(args)
	gain64 := fParse.Float64("g", 1.0,
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
}

//Original C version by Alexander Kritov 
//http://www.musicdsp.org/archive.php?classid=1#68  
func _DSF(x, a, N, fi float32) {
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
