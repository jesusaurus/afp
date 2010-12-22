package matrix

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

/**
 * Interleaves a 2d array into a 1d array
 * Allocates its own memory
 */
func Interleave(frame [][]float32) []float32 {
	var t int64 = 0
	out := make([]float32, len(frame) * len(frame[0]))
	
	for _,sample := range frame {
		for _,amplitude := range sample {
			out[t] = amplitude
			t++
		}
	}
	
	return out
}

/**
 * De-interleaves a 1d array into a 2d array
 */
func Deinterleave(frame []float32, samples int, channels int) [][]float32 {
	var t int64 = 0
	out := make([][]float32, samples)
	
	for i,_ := range(out) {
		out[i] = make([]float32, channels)
		for c,_ := range(out[i]) {
			out[i][c] = frame[t]
			t++
		}
	}
	
	return out
}
