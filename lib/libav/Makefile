# Copyright (c) 2010 Go Fightclub Authors
# This source code is released under the terms of the
# MIT license. Please see the file LICENSE for license details.

include $(GOROOT)/src/Make.inc

TARG=afp/libav
CGOFILES=libav.go
CGO_LDFLAGS=-lavformat -lavcodec -lz -lbz2
CGO_CFLAGS=-g
CLEANFILES+=main

libavtest: install main.go
	$(GC) main.go
	$(LD) -o $@ main.$O

include $(GOROOT)/src/Make.pkg

