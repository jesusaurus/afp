# Copyright (c) 2010 Go Fightclub Authors
# This source code is released under the terms of the
# MIT license. Please see the file LICENSE for license details.

include $(GOROOT)/src/Make.inc

TARG=afp/fftw
CGOFILES=fftw.go\
		constants.go
CGO_LDFLAGS=-lfftw -lfftwf -lm
CGO_CFLAGS=-g
#CLEANFILES+=main

include $(GOROOT)/src/Make.pkg

