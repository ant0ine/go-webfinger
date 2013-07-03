PREFIX=/usr/local/

all: webfinger

webfinger: main.go
	go build -v -o webfinger main.go

install:
	cp webfinger $(PREFIX)bin/
