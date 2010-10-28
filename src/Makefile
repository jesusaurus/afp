# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include /Users/eric/src/go/src/Make.inc

TARG=libav
CGOFILES=libav.go
CGO_LDFLAGS=-lavformat -lavcodec -lavutil -lavcore -lavdevice -lavfilter -lz -lbz2
CGO_CFLAGS=-g
CLEANFILES+=main 

libavtest: install main.go
			$(GC) main.go
			$(LD) -o $@ main.$O
			
include /Users/eric/src/go/src/Make.pkg

