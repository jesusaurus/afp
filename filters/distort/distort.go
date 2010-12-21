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
    amp Float32
    clipper func(*DistortFilter)
}

var clipTypes = map[string]func(*DistortFilter) {
    "hard" : hardCutoff,
}

func (self *DistortFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx

    fParse := flags.FlagParser(args)
    amp64 := fParse.Float64("a", math.NaN(),
        "The amplitude at which to clip the signal. 0.0 < a < 1.0")
    clipT := fParse.String("t", "hard", "The type of clipping used: hard")

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
