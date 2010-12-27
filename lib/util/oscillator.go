// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package util

import (
	"math"
)

const twoPi = 2 * math.Pi

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

func GetSawtootheOscillator(samplerate, freq, amp float32) (func() float32) {
	period := samplerate / freq //Roughly, period in slices
	delta := 2 * amp / period 
	var val float32 = 0

	return func() float32 {
		ret := val

		val += delta

		if val >= amp {
			val = -amp
		} 

		return ret
	}
}

func GetSquareOscillator(samplerate, freq, amp float32) (func() float32) {
	halfperiod := int(samplerate / (2 * freq)) //Roughly, period in slices
	val := amp
	i := 0
	
	return func() float32 {
		ret := val

		if i++; i >= halfperiod {
			i = 0
			val = -val
		}

		return ret
	}
}

func GetSineOscillator(samplerate, freq, amp float32) (func() float32) {
	period := samplerate / freq //Roughly, period in slices
	delta := twoPi / period
	var x float32 = 0

	return func() float32 {
		ret := amp * fastSine(x)

		x += delta

		if x >= twoPi {
			x -= twoPi
		} 

		return ret
	}
}

//fastSine evaluates the first 5 terms in the Taylor series:
//  sin(x) = x - x^3/3! + x^5/5! - x^7/7! + ...
//This provides Good Enough(TM) accuracy in the first two quadrants
//For [pi,2pi), we return -fastSine(x - pi)
//
//This function expects values only in the range [0,2pi)
//Outside this range, it will misbehave terribly
func fastSine(x float32) float32 {
	if x > math.Pi {
		return -fastSine(x - math.Pi)
	}

	x2 := x * x

	//We use a repeatedly factored representation to cut down on
	// FP mults, and allow us to do all our divisions at compile time
	return x * (1 - x2 * (1/6.0 - x2 * (1/120.0 - x2 * (1/5040.0 - x2 * 1/362880.0))))
}