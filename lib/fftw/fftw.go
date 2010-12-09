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

type fft_plan struct {
	plan C.fftwf_plan
	valid bool
}

type FFTPlan_r2r_1D_32 struct {
	fft_plan
	InBuff []float32
	OutBuff []float32
}

type FFTPlan_r2c_1D_32 struct {
	fft_plan
	InBuff []float32
	OutBuff []complex64
}

func (self *fft_plan) Execute() {
	if !self.valid {
		panic("Attempt to use destroyed FFT plan.")
	}
	C.fftwf_execute(self.plan)
}

func (self *fft_plan) Destroy() {
	self.valid = false
	C.fftwf_destroy_plan(self.plan)
}

//A note regarding the somewhat ugly cast from complex64 -> fftwf_complex
//fftwf_complex is a float[2], which has the same memory layout as Go's complex64
//The may seem brittle, but it's a standard leyout for complex numbers, and unlikely
//to change.
func NewRealToComplexPlan_1D_32(bufflen, flags int) *FFTPlan_r2c_1D_32 {
	inBuff := make([]float32, bufflen)
	outBuff := make([]complex64, bufflen / 2 + 1)

	return &FFTPlan_r2c_1D_32{fft_plan{C.fftwf_plan_dft_r2c_1d(C.int(bufflen), (*C.float)(&inBuff[0]), 
		(*C.fftwf_complex)(unsafe.Pointer(&outBuff[0])),
				C.uint(C.FFTW_UNALIGNED | flags)), true}, inBuff, outBuff}
}

func NewFFTPlan_r2r_1D_32(bufflen, flags int, kind C.fftwf_r2r_kind) *FFTPlan_r2r_1D_32 {
	inBuff := make([]float32, bufflen)
	outBuff := make([]float32, bufflen)

	return &FFTPlan_r2r_1D_32{fft_plan{C.fftwf_plan_r2r_1d(C.int(bufflen), (*C.float)(&inBuff[0]),
				(*C.float)(&inBuff[0]), kind, C.uint(C.FFTW_UNALIGNED | flags)), true}, inBuff, outBuff}

}

/*
func NewRealToRealPlan_2D_32(data, output []float32, flags int, kind C.fftwf_r2r_kind) *FFTPlan32 {
     
	fftwf_plan_r2r_2d(int n0, int n1, double *in, double *out,
        fftw_r2r_kind kind0, fftw_r2r_kind kind1,unsigned flags)
}
*/