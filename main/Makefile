GV=6

test: 
	$(GV)g -o main.$(GV) testmain.go filterlist.go startup.go
	$(GV)l -o afp main.$(GV)

clean:
	-rm main.$(GV) *~ afp 2> /dev/null