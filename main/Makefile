# Copyright (c) 2010 Go Fightclub Authors
# This source code is released under the terms of the
# MIT license. Please see the file LICENSE for license details.

include $(GOROOT)/src/Make.inc

afp: main.go filterlist.go startup.go
	$(GC) -o main.$(O) main.go filterlist.go startup.go pipeparse.go
	$(LD) -o ../afp main.$(O)

test: testmain.go filterlist.go startup.go
	$(GC) -o main.$(O) testmain.go filterlist.go startup.go pipeparse.go
	$(LD)l -o ../afptest main.$(O)

clean:
	-rm main.$(O) *~ ../afp 2> /dev/null