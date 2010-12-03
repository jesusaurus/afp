include $(GOROOT)/src/Make.inc
TARG=fightclub
CGOFILES=readin.go readout.go
CGO_CFLAGS=`pkg-config --cflags libavcodec libavformat`
CGO_LDFLAGS=`pkg-config --libs libavcodec libavformat`
include $(GOROOT)/src/Make.pkg 

CLEANFILES+=open write $(PKGDIR)/$(TARG).a 

open: install open.go
	$(GC) open.go
	$(LD) -o $@ open.$O

write: install write.go
	$(GC) write.go
	$(LD) -o $@ write.$O

