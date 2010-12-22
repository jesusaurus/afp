// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package util

/**
 * Given a slice of samples, returns a slice of channels, or visa-versa
 */
func Invert(frame [][]float32) [][]float32 {
	out := make([][]float32, len(frame[0]))

	for i := range out {
		out[i] = make([]float32, len(frame))
	}

	for i := range frame {
		for j := range frame[i] {
			out[j][i] = frame[i][j]
		}
	}

	return out
}


/**
 * Extracts one channel of audio data into a contiguous slice
 */
func ExtractChannel(frame [][]float32, channel int) []float32 {
	out := make([]float32, len(frame))

	if channel > len(frame[0]) {
		panic("Attempt to extract a non-existent channel")
	}

	for i := range frame {
		out[i] = frame[i][channel]
	}

	return out
}
