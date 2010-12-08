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

type FFTPlan32 struct {
	plan C.fftwf_plan
}

func NewRealToComplexPlan_1D_32(data []float32, output []complex64, flags int) *FFTPlan32 {
	return &FFTPlan32{C.fftwf_plan_dft_r2c_1d(C.int(len(data)), (*C.float)(&data[0]), 
		(*C.fftwf_complex)(unsafe.Pointer(&output[0])),
			C.uint(C.FFTW_UNALIGNED | flags))}
}

func NewRealToRealPlan_1D_32(data, output []float32, flags int, kind C.fftwf_r2r_kind) *FFTPlan32 {

	return &FFTPlan32{C.fftwf_plan_r2r_1d(C.int(len(data)), (*C.float)(&data[0]),
			           (*C.float)(&output[0]), kind, C.uint(C.FFTW_UNALIGNED | flags))}
}

func (self *FFTPlan32) Execute() {
	C.fftwf_execute(self.plan)
}
/*
     fftw_plan fftw_plan_r2r_2d(int n0, int n1, double *in, double *out,
                                fftw_r2r_kind kind0, fftw_r2r_kind kind1,
    unsigned flags)
*/