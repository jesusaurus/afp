package tonefilter

import (
	"afp"
	"afp/flags"
	"math"
	"os"
	"fmt"
)

const (
	TONE_SINE = iota
	TONE_SAW
	TONE_SQUARE
)

type ToneSource struct {
	context *afp.Context
	header afp.StreamHeader

	toneType int
	toneFrequency float32
	toneAmplitude float32
	toneLength float32
	
	buffers [][][]float32
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
	var err os.Error;
	
	self.context = ctx
	
	parser := flags.FlagParser(args)
	var toneType *string = parser.String("t", "sine", "Type of tone to generate: sine, saw, square (sine)")
	var freq *float = parser.Float("f", 440.0, "The frequency of the tone to generate (440Hz)")
	var amp *float = parser.Float("a", 0.75, "The amplitude of the output tone (0.75)")
	var len *float = parser.Float("l", 10.0, "The length of the tone in seconds (10s)")
	var channels *int = parser.Int("c", 2, "The number of channels to generate (2)")
	var rate *int = parser.Int("r", 44100, "The sampling rate of the output (44.1khz)")
	parser.Parse()
	
	dumpConfig(*toneType, *freq, *amp, *len, *channels, *rate)
	
	self.toneType, err = mapToneType(*toneType)
	self.toneFrequency = float32(*freq)
	self.toneAmplitude = float32(*amp)
	self.toneLength = float32(*len)
	
	self.header.Version = 1
	self.header.Channels = int8(*channels)
	self.header.SampleSize = 4
	self.header.SampleRate = int32(*rate)
	self.header.FrameSize = 4096
	self.header.ContentLength = 0

	self.makeBuffers(2, self.header.FrameSize, self.header.Channels)	
	
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
	var toneMap map[string]int = map[string]int {
		"sine"		: TONE_SINE,
		"saw"		: TONE_SAW,
		"square"	: TONE_SQUARE,
	}
	return toneMap[toneName], nil
}

func (self *ToneSource) Start() {
	self.context.HeaderSink <- self.header;
	
	var (
		t, dt float32
		fo int32 = 0		/* frame offset */
		cb int = 0			/* current buffer */
		c int8
	)
	
	dt = 1.0 / float32(self.header.SampleRate)	
	
	for t = 0; t < self.toneLength; t += dt {
		for c = 0; c < self.header.Channels; c++ {
			self.buffers[cb][fo][c] = float32(math.Sin(float64(t * 2 * math.Pi * self.toneFrequency))) * self.toneAmplitude
		}
		fo++
		
		if fo == self.header.FrameSize {
			self.context.Sink <- self.buffers[cb]
			cb = 1 - cb
			fo = 0
		}
	}
}

func (self *ToneSource) makeBuffers(numBuffers int, frameSize int32, channels int8) {
	/* initialize floatSamples buffer */
	self.buffers = make([][][]float32, numBuffers)
	for i,_ := range(self.buffers) {
		self.buffers[i] = make([][]float32, frameSize)

		for j,_ := range(self.buffers[i]) {
			self.buffers[i][j] = make([]float32, channels)
		}
	}
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