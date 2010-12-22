package main

import (
	"afp/fftw"
)

func main() {
	for i := 0; i < 1000; i++ {
		d := make([]float32, 1152)
		_ = fftw.RealToReal1D_32(d, false, 3, fftw.MEASURE, fftw.R2HC)
		_ = fftw.RealToReal1D_32(d, false, 3, fftw.MEASURE, fftw.HC2R)
		print(".")
	}
}
