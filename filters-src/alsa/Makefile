include $(GOROOT)/src/Make.inc
TARG=alsa
CGOFILES=alsa.go
CGO_CFLAGS=`pkg-config --cflags alsa`
CGO_LDFLAGS=`pkg-config --libs alsa`
include $(GOROOT)/src/Make.pkg

CLEANFILES+=alsa $(PKGDIR)/$(TARG).a
