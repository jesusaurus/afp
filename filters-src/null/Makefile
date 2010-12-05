include $(GOROOT)/src/Make.inc

PKG_NAME=null
GOFILES=\
	null.go

all:
	$(GC) -o $(PKG_NAME).$(O) $(GOFILES)

install:
	cp $(PKG_NAME).$(O) $(PROJ_ROOT)/filters