package mono

import "afp"

type MonoFilter struct {
	ctx *afp.Context
}

func (self *MonoFilter) GetType() int {
	return afp.PIPE_LINK
}

func (self *MonoFilter) Init(ctx *afp.Context, args []string) os.Error {
	self.ctx = ctx
	parser := flags.FlagParser(args)
	parser.Parse()
	if len(parser.Args()) {
		return os.NewError("mono link takes 0 arguments")
	}
	return nil
}

func (self *MonoFilter) Start() {
	header StreamHeader
	ctx := self.ctx
	header <-ctx.HeaderSource
	// Unpack header
	channels := header.Channels

	header.Channels = 1
	// TODO: Is this math correct? Is this guaranteed to be accurate?
	header.ContentLength = header.ContentLength / channels
	ctx.HeaderSink <-header

	if channels == 1 { // Already mono, don't manipulate
		for frame := ctx.Source {
			ctx.Sink <- frame
		}
	} else {
		for frame := ctx.Source {
			ctx.Sink(mergeChannels(frame, channels))
		}
	}
}

// Merge channels, return new frame. May modify the frame array.
func mergeChannels(frame [][]float32, channels int) [][]float32


// mergeChannels for frame[channel #][sample]
// This _might_ be significantly faster than the other versions
// Optimization: Merge all channels into channel one instead of
// allocating a new array.
func mergeChannels(frame [][]float32, channels int) [][]float32 {
	var monoValue float32
	// Use sample indexes from the first channel
	for sample, _ := range frame[0] {
		monoValue = 0.0
		for _, channel := range frame {
			monoValue += channel[sample] / channels
		}
		// Assign to first channel
		frame[0][sample] := monoValue
	}

	newFrame = make([][]float32, 1);
	newFrame[0] = frame[0]
	return newFrame
}


// mergeChannels for frame[sample][channel #]
// Optimization:
func mergeChannels(frame [][]float32, channels int) [][]float32 {
	var monoValue float32
	for sample, sampleValues := range frame {
		monoValue = 0.0
		for _, channelValue := range sampleValues {
			monoValue += channelValue / channels
		}

		// Allocate a new length one array for each sample: performance
		// issue?
		newSample = make([]float32, 1)
		newSample[0] = monoValue
		frame[sample] = newSample
	}
	return frame
}


// merge Channels for frame[sample][channel #]
// Uses a slicing optimization to cut down on allocations at the
// expense of leaving many open spaces and some cleverness
func mergeChannels(frame [][]float32, channels int) [][]float32 {
	var monoValue float32
	for sample, sampleValues := range frame {
		for _, channelValue := range sampleValues {
			monoValue += sample / channels
		}
		// Possible trade off: memory for allocation time, *cleverness*
		frame[sample][0] = monoValue
		frame[sample] = frame[sample][:1] // assign 1 length slice for mono channel
	}
	return frame
}
