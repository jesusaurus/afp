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

// mergeChannels for frame[channel #][sample]
// This _might_ be significantly faster than the other versions
func mergeChannels(frame [][]float32, channels int) [][]float32 {
	for slice, _ := range frame[0] {
		for channel, _ := range frame {
			mval += frame[channel][slice] / channels
		}
		// Reuse the first channel
		frame[0][slice] := mval
	}
	// Return a slice with just the first channel
	return frame[:1]
}

// mergeChannels for frame[sample][channel #]
func mergeChannels(frame [][]float32, channels int) [][]float32 {
	for i, samples := range frame {
		for j, sample := range samples {
			mval += sample / channels
		}
		// Allocate a new length one array for each one: time issue?
		sampleLine = make([]float32, 1)
		sampleLine[0] = mval
		frame[i] = mval
	}
	return frame
}

// merge Channels for frame[sample][channel #]
// Uses a slicing optimization to cut down on allocations
func mergeChannels(frame [][]float32, channels int) [][]float32 {
	for i, samples := range frame {
		for j, sample := range samples {
			mval += sample / channels
		}
		// Assign to the first line item, assign a slice
		// Possible trade off: memory for allocation time.
		frame[i][0] = mval
		frame[i] = samples[:1]
	}
	return frame
}
