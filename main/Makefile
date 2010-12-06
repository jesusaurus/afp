# Copyright (c) 2010 Go Fightclub Authors
# This source code is released under the terms of the
# MIT license. Please see the file LICENSE for license details.

GV=6

test: testmain.go filterlist.go startup.go
	$(GV)g -o main.$(GV) testmain.go filterlist.go startup.go
	$(GV)l -o ../afptest main.$(GV)

afp: main.go filterlist.go startup.go
	$(GV)g -o main.$(GV) main.go filterlist.go startup.go
	$(GV)l -o ../afp main.$(GV)

clean:
	-rm main.$(GV) *~ ../afp 2> /dev/null