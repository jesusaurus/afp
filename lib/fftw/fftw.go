// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

/*
 This is a set of wrappers around the fftw3 library.  For more information
 and use documentation visit http://www.fftw.org
*/

package fftw

//#include <fftw3.h>
import "C"

import (
	"unsafe"
	)

//A note regarding the somewhat ugly cast from complex64 -> fftwf_complex
//fftwf_complex is a float[2], which has the same memory layout as Go's complex64
//The may seem brittle, but it's a standard leyout for complex numbers, and unlikely
//to change.
func RealToComplex1D_32(data []float32, iterations, flags int) []complex64 {
	output := make([]complex64, len(data) / 2 + 1)
	
	plan := C.fftwf_plan_dft_r2c_1d(C.int(len(data)), (*C.float)(&data[0]), 
		(*C.fftwf_complex)(unsafe.Pointer(&output[0])),
		C.uint(C.FFTW_UNALIGNED | flags))

	for i := 0; i < iterations; i++ {
		C.fftwf_execute(plan)
	}

	return output
}

func RealToReal1D_32(data []float32, inPlace bool, iterations, flags int, kind C.fftwf_r2r_kind) (out []float32) {
	if inPlace {
		out = data
	} else {
		out = make([]float32, len(data))
	}

	plan :=  C.fftwf_plan_r2r_1d(C.int(len(data)), (*C.float)(&data[0]), (*C.float)(&out[0]),
		kind, C.uint(C.FFTW_UNALIGNED | flags))

	for i := 0; i < iterations; i++ {
		C.fftwf_execute(plan)
	}

	return

}
/*
     fftw_plan fftw_plan_r2r_2d(int n0, int n1, double *in, double *out,
                                fftw_r2r_kind kind0, fftw_r2r_kind kind1,
    unsigned flags)
*/
