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

	C.fftwf_destroy_plan(plan)
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

	C.fftwf_destroy_plan(plan)
	return
}

func RealToReal2D_32(data [][]float32, iterations, flags int, 
	kind0 C.fftwf_r2r_kind, kind1 C.fftwf_r2r_kind) (out [][]float32) {

	//fftw expects a 2d array in row major order.
	//We have no guarantees about the memory layout 
	//of any frame, so we have to copy and flatten
	flattenedData := make([]float32, len(data) * len(data[0]))
	
	p := 0
	for i := range data {
		for j := range data[i] {
			flattenedData[p] = data[i][j]
			p++
		}
	}

	flattenedOut := make([]float32, len(data) * len(data[0]))
	
	plan :=  C.fftwf_plan_r2r_2d(C.int(len(data)), C.int(len(data[0])), 
		(*C.float)(&flattenedData[0]), (*C.float)(&flattenedOut[0]),
		kind0, kind1, C.uint(C.FFTW_UNALIGNED | flags))

	for i := 0; i < iterations; i++ {
		C.fftwf_execute(plan)
	}


	//Slice the flattened array for return to user
	out = make([][]float32, len(data))

	sampleLen := len(data[0])
	slice := 0
	for i := range out {
		out[i] = flattenedOut[slice:slice + sampleLen]
		slice += sampleLen
	} 

	C.fftwf_destroy_plan(plan)
	return
}
