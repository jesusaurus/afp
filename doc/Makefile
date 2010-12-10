ifeq ($(AFPMANDIR),"")
$(error $$AFPMANDIR is not set)
endif

install: 
	-mkdir $(AFPMANDIR)/man1/ $(AFPMANDIR)/man3/ $(AFPMANDIR)/man5/ 2> /dev/null	
	cp afp.1 $(AFPMANDIR)/man1/
	cp afp-dev.3 $(AFPMANDIR)/man3/
	cp afp-fmt.5 $(AFPMANDIR)/man5/
	mandb -c