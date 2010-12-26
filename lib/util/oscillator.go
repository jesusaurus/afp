// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package util

import (
	"afp"
	)

func GetTriangleOscillator(samplerate, freq, amp float32) (func() float32) {
	period := samplerate / freq //Roughly, period in slices
	delta := 4 * amp / period 
	var val float32 = 0

	return func() float32 {
		ret := val

		val += delta

		if val >= amp || val <= -amp {
			delta = -delta
		} 

		return ret
	}
}

