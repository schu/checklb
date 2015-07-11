version = $(shell git describe --tags --always --dirty)

all: build

build:
	godep restore
	go build -ldflags "-X main.Version '$(version)'" checklb.go
