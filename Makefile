include $(GOROOT)/src/Make.inc

TARG=main
GOFILES=\
	main.go\

include $(GOROOT)/src/Make.pkg

link:
	6l -o main _go_.6
