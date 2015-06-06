version = $(shell git describe --tags --always)

all: build

build:
	godep restore
	go build -ldflags "-X main.Version '$(version)'" checklb.go
