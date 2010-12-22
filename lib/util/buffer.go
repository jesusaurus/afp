// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package util

import "afp"


//Buffer will buffer n frames except in the case that the stream 
//is closed before n frames have been read.
//For performance reasons, This should be used only where n >> CHAN_BUF_LEN
//due to the necessity of buffering samples in a second channel. 
//For n <= CHAN_BUF_LEN Buffer will delegate to FastBuffer.
func Buffer(n int, source <-chan [][]float32) <-chan [][]float32 {
	if n <= afp.CHAN_BUF_LEN {
		FastBuffer(frames, source)
		return source;
	}

	buff := make(chan [][]float32, n)
	buffered := 0

	for s := range source {
		buff <- s

		if buffered >= n {
			break
		}
	}
	
	//We need to copy all subsequent frames
	//sent into the chan the caller will
	//now be reading from.
	go func() {
		for s := range source {
			buff <- s
		}
	}()

	return buff
}

func FastBuffer(frames int) int {
}