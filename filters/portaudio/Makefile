#Copyright (c) 2010 AFP Authors
#This source code is released under the terms of the
#MIT license. Please see the file LICENSE for license details.

include $(GOROOT)/src/Make.inc

TARG=afp/filters/portaudio
CGOFILES=portaudio.go
CGO_CFLAGS=`pkg-config --cflags portaudio-2.0` -g
CGO_LDFLAGS=`pkg-config --libs portaudio-2.0` -lportaudio -lpthread

include $(GOROOT)/src/Make.pkg

CLEANFILES+=portaudio $(pkgdir)/$(TARG).a
