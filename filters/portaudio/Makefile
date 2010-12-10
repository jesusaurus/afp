include $(GOROOT)/src/Make.inc

TARG=afp/filters/portaudio
CGOFILES=portaudio.go
CGO_CFLAGS=`pkg-config --cflags portaudio-2.0`
CGO_LDFLAGS=`pkg-config --libs portaudio-2.0` -lpthread

include $(GOROOT)/src/Make.pkg

CLEANFILES+=portaudio $(PKGDIR)/$(TARG).a
