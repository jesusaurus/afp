include $(GOROOT)/src/Make.inc

PKG_NAME=fexec
GOFILES=\
		exec.go\
		execSource.go\
		execSink.go\
		execLink.go

build:
	$(GC) -o $(PKG_NAME).$(O) $(GOFILES)

install:
	cp $(PKG_NAME).$(O) $(PROJ_ROOT)/filters
