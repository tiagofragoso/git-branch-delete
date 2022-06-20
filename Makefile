BINDIR=bin
NAME=git-branch-delete

build: 
	go build -o $(BINDIR)/$(NAME)

clean: 
	rm $(BINDIR)/*
