# Copyright (c) 2010 AFP Authors
# This source code is released under the terms of the
# MIT license. Please see the file LICENSE for license details.

include $(GOROOT)/src/Make.inc

TARG=afp/fftw
CGOFILES=fftw.go\
		constants.go
CGO_LDFLAGS=-lfftw3 -lfftw3f -lm
CGO_CFLAGS=-g
#CLEANFILES+=main

include $(GOROOT)/src/Make.pkg

bench_fftw: bench_fftw.go
	$(GC) bench_fftw.go
	$(LD) -o bench_fftw bench_fftw.$(O)
