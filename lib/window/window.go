// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package window

import (
	"math"
)

func Hann(samples []float32) []float32 {
	n := len(samples)
	windowed := make([]float32, n)
	
	for i, amp := range(samples) {
		windowed[i] = amp * float32(0.5 * (1.0 - math.Cos(2 * math.Pi * float64(i) / float64(n - 1))))
	}
	
	return windowed
}

