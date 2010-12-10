include $(GOROOT)/src/Make.inc

TARG=afp/filters/portaudio
CGOFILES=portaudio.go
CGO_CFLAGS=
CGO_LDFLAGS=-lportaudio -lpthread

include $(GOROOT)/src/Make.pkg

CLEANFILES+=portaudio $(pkgdir)/$(TARG).a
