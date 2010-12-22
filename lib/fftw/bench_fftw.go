package main

import (
	"afp/fftw"
)

func main() {
	for i := 0; i < 1000; i++ {
		d := make([]float32, 5192)
		_ = fftw.RealToReal1D_32(d, false, 3, fftw.EXHAUSTIVE, fftw.DHT)
		print(".")
	}
}
