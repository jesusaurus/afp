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
}

func (self *MonoFilter) Start() {
    header StreamHeader
    ctx := self.ctx
    header <- ctx.HeaderSource
    // Unpack header
    channels := header.Channels

    // Modify the header, then send it to the sink
    header.Channels = 1
    // TODO: Is this math correct?
    header.ContentLength = header.ContentLength / channels
    ctx.HeaderSink <- header

    if channels == 0 {
		for audio := ctx.Source {
			ctx.Sink <- audio
		}
	} else {
	    for audio := ctx.Source {
			ctx.Sink(mergeChannels(audio, channels))
		}
	}
}

// mergeChannels for audio[channel #][sample]
func mergeChannels(audio [][]float32, channels int) [][]float32 {
	for i, _ := range audio[0] {
		for channel, _ := range audio {
			mval += audio[channel][i] / channels
		}
		// Reuse the first channel
		audio[0][i] := mval
	}
	// Return a slice with just the first channel
	return audio[:1]
}

// mergeChannels for audio[sample][channel #]
func mergeChannels(audio [][]float32, channels int) [][]float32 {
	for i, samples := range audio {
		for j, sample := range samples {
			mval += sample / channels
		}
		sampleLine = make([]float32, 1)
		sampleLine[0] = mval
		audio[i] = mval
	}
	return audio
}
