BINDIR			:= bin
LDFLAGS			:= -w -s

.PHONY: all
all: build

.PHONY: build
build:
	go get
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
