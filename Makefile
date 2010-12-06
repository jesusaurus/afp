build:
	./build.sh
install:
	ifndef $(AFPINSTALLDIR)
	$(error $$AFPINSTALLDIR is not set)
	else
	cp afp $(AFPINSTALLDIR)
	endif

	ifndef $(AFPMANDIR)
	@echo warning: $$AFPMANDIR is not set, no manpages will be installed
	else 
	cd doc && make install
	endif

clean:
	./clean.sh