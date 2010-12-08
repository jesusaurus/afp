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

	//Real to Real transform kinds - documentation from 
	//http://www.fftw.org/fftw3_doc/Real_002dto_002dReal-Transform-Kinds.html#Real_002dto_002dReal-Transform-Kinds
	
	//computes a real-input DFT with output in “halfcomplex” format, i.e. 
	//real and imaginary parts for a transform of size n stored as:
	//r0, r1, r2, ..., rn/2, i(n+1)/2-1, ..., i2, i1 (Logical N=n, inverse is FFTW_HC2R.)
	R2HC = C.FFTW_R2HC 

	//Computes the reverse of FFTW_R2HC, above. (Logical N=n, inverse is FFTW_R2HC.)
	HC2R = C.FFTW_HC2R 
	
	//Computes a discrete Hartley transform. (Logical N=n, inverse is FFTW_DHT.)
	DHT = C.FFTW_DHT

	//Computes an REDFT00 transform, i.e. a DCT-I. (Logical N=2*(n-1), inverse is FFTW_REDFT00.) 
	REDFT00 = C.FFTW_REDFT00
	
	//Computes an REDFT10 transform, i.e. a DCT-II (sometimes called “the” DCT). 
	//(Logical N=2*n, inverse is FFTW_REDFT01.)
	REDFT10 = C.FFTW_REDFT10 
	
	//Computes an REDFT01 transform, i.e. a DCT-III (sometimes called “the” IDCT,
	// being the inverse of DCT-II). (Logical N=2*n, inverse is FFTW_REDFT=10.)
	REDFT01 = C.FFTW_REDFT01 
	
	//Computes an REDFT11 transform, i.e. a DCT-IV. (Logical N=2*n, inverse is FFTW_REDFT11.)
	REDFT11 = C.FFTW_REDFT11 
	
	//Computes an RODFT00 transform, i.e. a DST-I. (Logical N=2*(n+1), inverse is FFTW_RODFT00.)
	RODFT00 = C.FFTW_RODFT00

	//Computes an RODFT10 transform, i.e. a DST-II. (Logical N=2*n, inverse is FFTW_RODFT01.)
	RODFT10 = C.FFTW_RODFT10

	//Computes an RODFT01 transform, i.e. a DST-III. (Logical N=2*n, inverse is FFTW_RODFT=10.)
	RODFT01 = C.FFTW_RODFT01

	//Computes an RODFT11 transform, i.e. a DST-IV. (Logical N=2*n, inverse is FFTW_RODFT11.)
	RODFT11 = C.FFTW_RODFT11
	)

func RealToComplex1D_32(data []float32, iterations, flags int) []complex64 {
	output := make([]complex64, len(data) / 2 + 1)
	plan := C.fftwf_plan_dft_r2c_1d(len(data), &data[0], &output[0],
		                            C.FFTW_UNALIGNED | flags)

	for i := 0; i < iterations; i++ {
		C.fftwf_execute(plan)
	}

	return output;
}

func RealToReal1D_32(data []float32, iterations, flags int) []float32 {
    plan :=  C.fftwf_plan_r2r_1d(len(data), &data[0], &data[0],
		fftw_r2r_kind kind, C.FFTW_UNALIGNED | flags);


}
/*
     fftw_plan fftw_plan_r2r_1d(int n, double *in, double *out,
    fftw_r2r_kind kind, unsigned flags);
     fftw_plan fftw_plan_r2r_2d(int n0, int n1, double *in, double *out,
                                fftw_r2r_kind kind0, fftw_r2r_kind kind1,
    unsigned flags)
*/