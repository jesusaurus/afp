#Copyright (c) 2010 Go Fightclub Authors
#This source code is released under the terms of the
#MIT license. Please see the file LICENSE for license details.

include $(GOROOT)/src/Make.inc

TARG=afp/filters/alsa
CGOFILES=alsa.go
CGO_CFLAGS=`pkg-config --cflags alsa`
CGO_LDFLAGS=`pkg-config --libs alsa`
include $(GOROOT)/src/Make.pkg

CLEANFILES+=alsa $(PKGDIR)/$(TARG).a
