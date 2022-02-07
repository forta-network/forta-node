prefix = /usr/local

all:

install:
	wget -O forta ARTIFACTS_URL/forta-REVISION
	install -D forta \
		$(DESTDIR)$(prefix)/bin/forta

clean:
	-rm -f forta

distclean: clean

uninstall:
	-rm -f $(DESTDIR)$(prefix)/bin/forta

.PHONY: all install clean distclean uninstall
