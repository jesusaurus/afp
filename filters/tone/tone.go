// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.


package tonefilter


import (
	"afp"
	"afp/flags"
	"math"
	"os"
	"fmt"
)


const (
	FRAME_SIZE = 4096
	TONE_SINE  = iota
	TONE_SAW
	TONE_SQUARE
)


type ToneSource struct {
	context *afp.Context
	header  afp.StreamHeader

	toneType      int
	toneFrequency float32
	toneAmplitude float32
	toneLength    float32
}


/**
 * initialize the filter storage
 */

func NewToneSource() afp.Filter {
	return &ToneSource{}
}

/**
 * ToneSource is unsurprisingly a source
 */

func (self *ToneSource) GetType() int {
	return afp.PIPE_SOURCE
}

/**
 * Configure settings with hopefully reasonable defaults
 */

func (self *ToneSource) Init(ctx *afp.Context, args []string) os.Error {
	var err os.Error

	self.context = ctx

	parser := flags.FlagParser(args)
	var toneType *string = parser.String("t", "sine", "Type of tone to generate: sine, saw, square (sine)")
	var freq *float = parser.Float("f", 440.0, "The frequency of the tone to generate")
	var amp *float = parser.Float("a", 0.75, "The amplitude of the output tone")
	var len *float = parser.Float("l", 10.0, "The length of the tone in seconds")
	var channels *int = parser.Int("c", 2, "The number of channels to generate")
	var rate *int = parser.Int("r", 44100, "The sampling rate of the output")
	parser.Parse()

	/*	dumpConfig(*toneType, *freq, *amp, *len, *channels, *rate)*/

	self.toneType, err = mapToneType(*toneType)
	self.toneFrequency = float32(*freq)
	self.toneAmplitude = float32(*amp)
	self.toneLength = float32(*len)

	self.header.Version = 1
	self.header.Channels = int8(*channels)
	self.header.SampleSize = 4
	self.header.SampleRate = int32(*rate)
	self.header.FrameSize = FRAME_SIZE
	self.header.ContentLength = 0

	return err
}

func dumpConfig(toneType string, freq float, amp float, len float, channels int, rate int) {
	fmt.Fprintf(os.Stderr, "Tone Config (%s):", toneType)
	fmt.Fprintf(os.Stderr, "\n  frequency: %f", freq)
	fmt.Fprintf(os.Stderr, "\n  amplitude: %f", amp)
	fmt.Fprintf(os.Stderr, "\n  length: %f", len)
	fmt.Fprintf(os.Stderr, "\n  channels: %d", channels)
	fmt.Fprintf(os.Stderr, "\n  sample rate: %d\n\n", rate)
}

/**
 * Given a string, what is the corresponding tone type?
 */

func mapToneType(toneName string) (int, os.Error) {
	var toneMap map[string]int = map[string]int{
		"sine":   TONE_SINE,
		"saw":    TONE_SAW,
		"square": TONE_SQUARE,
	}
	return toneMap[toneName], nil
}

func (self *ToneSource) Start() {
	self.context.HeaderSink <- self.header

	var (
		s      int64            /* sample */
		t      float64          /* time */
		fo     int32        = 0 /* frame offset */
		c      int8             /* channel iterator */
		buffer *[][]float32 = makeBuffer(self.header.FrameSize, self.header.Channels)
	)

	for s = 0; s < int64(self.toneLength*float32(self.header.SampleRate)); s++ {
		t = float64(s) / float64(self.header.SampleRate)

		for c = 0; c < self.header.Channels; c++ {
			(*buffer)[fo][c] = float32(math.Sin(t*2.0*math.Pi*float64(self.toneFrequency))) * self.toneAmplitude
		}
		fo++

		if fo == self.header.FrameSize {
			self.context.Sink <- *buffer

			buffer = makeBuffer(self.header.FrameSize, self.header.Channels)
			fo = 0
		}
	}

	if fo != self.header.FrameSize {
		fo += 1
		for fo < self.header.FrameSize {
			for c = 0; c < self.header.Channels; c++ {
				(*buffer)[fo][c] = 0.0
			}
			fo += 1
		}
		self.context.Sink <- (*buffer)
	}
}

func makeBuffer(size int32, channels int8) *[][]float32 {
	b := make([][]float32, size)
	for i, _ := range b {
		b[i] = make([]float32, channels)
	}

	return &b
}

func (self *ToneSource) Stop() os.Error {
	close(self.context.Sink)
	return nil
}

/*type Filter interface {
	GetType() int
	Init(*Context, []string) os.Error
	Start()
	Stop() os.Error
}
*/
