include $(GOROOT)/src/Make.$(GOARCH)

TARG=notifo
GOFILES=\
	notifo.go\

include $(GOROOT)/src/Make.pkg
