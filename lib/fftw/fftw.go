// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

/*
 This is a set of wrappers around the fftw3 library.  For more information
 and use documentation visit http://www.fftw.org
*/

package fftw

/*
#include <fftw3.h>
 

*/
import "C"

import (
	"unsafe"
	)

const (
	//Plan building constants from fastest to best
	//guarantee of optimality
	ESTIMATE = C.FFTW_ESTIMATE
	MEASURE = C.FFTW_MEASURE //The default
	PATIENT = C.FFTW_PATIENT
	EXHAUSTIVE = C.FFTW_EXHAUSTIVE
	WISDOM_ONLY = C.FFTW_WISDOM_ONLY

	//Algorithm Restriction flags
	DESTROY_INPUT = C.FFTW_DESTROY_INPUT //Applies only to out of place transforms
	PRESERVE_INPUT = C.FFTW_PRESERVE_INPUT
	//FFTW_UNALIGNED is always passed, since arrays are allocated in Go code
	//And thus have no alignment guarantees known to fftw
	)



func RealToReal1D_32(data []float32, iterations, flags int) []float32 {
	C.fftw_plan

	for i := 0; i < iterations; i++ {
	}
}

func RealToReal1D_32(data []float32, iterations, flags int) []float32 {
}

func RealToComplex1D_32(data []float32, iterations, flags int) []complex64 {
	panic("Unimplemented")
}

func RealToComplex1D_64(data []float64, iterations, flags int) []complex128 {
	panic("Unimplemented")
}

func ComplexToReal1D_32(data []complex64, iterations, flags int) []float32 {
	panic("Unimplemented")
}

func ComplexToReal1D_32(data []complex128, iterations, flags int) []float64 {
	panic("Unimplemented")
}
