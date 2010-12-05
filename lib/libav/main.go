// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main;

import (
	"libav"
	"unsafe"
	"os"
	"flag"
	"fmt"
	"encoding/binary" 
)

const (
	samplesPerSecond = 44100
	samplesPerMillisecond = samplesPerSecond / 1000
	delayTimeInMs = 605
	channels = 2
	bytesPerSample = 2
)

func main() {
	var context libav.AVDecodeContext
	var currBuffer = 0
	var buffer [2][]int16
	var extraSamples int = samplesPerMillisecond * delayTimeInMs * channels
	var infile string

	buffer[0] = make([]int16, extraSamples, extraSamples)
	
	libav.InitDecoding()
	
	flag.Parse() // Scans the arg list and sets up flags
	if flag.NArg() > 0 {
	    infile = flag.Arg(0)
	} else {
		panic("Dumb!")
	}

	libav.PrepareDecoding(infile, &context)
	
	var (
		frame = 0
		l = 1
		length = 0
	)
	
	for l > 0 {
		l = libav.DecodePacket(context)
		
		if l > 0 {
			numberOfSamples := l / bytesPerSample
			fmt.Fprintf(os.Stderr, "Frame: %d Length: %d Number of samples: %d\n", frame, l, numberOfSamples)
			frame += 1
			length += numberOfSamples
			decodedSamples := (*(*[1 << 31 - 1]int16)(unsafe.Pointer(context.Context.Outbuf)))[:numberOfSamples]

/*			os.Stdout.Write((*(*[1 << 31 - 1]uint8)(unsafe.Pointer(context.Context.Outbuf)))[:(numberOfSamples*bytesPerSample)])*/
			buffer[1 - currBuffer] = append(buffer[currBuffer], decodedSamples...)
			currBuffer = 1 - currBuffer;
		} else {
			fmt.Fprintf(os.Stderr, "Frame: %d Length: %d\n", frame, l)
		}
	}
	
	buffer[1 - currBuffer] = make([]int16, len(buffer[currBuffer]) + extraSamples, len(buffer[currBuffer]) + extraSamples)
	for t0,_ := range buffer[1 - currBuffer] {
		if t0 > extraSamples {
			if t0 < len(buffer[currBuffer]) {
				buffer[1 - currBuffer][t0] = int16(0.8 * (float32(buffer[currBuffer][t0]) + (float32(buffer[currBuffer][t0 - extraSamples]) * 0.4)))
			} else {
				buffer[1 - currBuffer][t0] = int16(float32(buffer[currBuffer][t0 - extraSamples]) * 0.4)
			}
		}
	}

	currBuffer = 1 - currBuffer;

/*	fmt.Fprintf(os.Stderr, "Len: %d Buffer 0: %d Buffer 1: %d\n", length, len(buffer[0]), len(buffer [1]))*/

/*	bufferBytes := (*(*[1 << 31 - 1]uint8)(unsafe.Pointer(&buffer[currBuffer])))[:(len(buffer[currBuffer]) * 2)]
	bytes, err := os.Stdout.Write(bufferBytes) */
/*	bufferBytes := (*(*[]uint8)(unsafe.Pointer(&buffer[currBuffer])))[:(2*len(buffer[currBuffer]))]*/
	err := binary.Write(os.Stdout, binary.LittleEndian, buffer[currBuffer]) 
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
	}
}
