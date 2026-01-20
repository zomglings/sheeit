VERSION := $(shell grep 'var Version' version/version.go | cut -d'"' -f2)

.PHONY: build clean

build:
	@mkdir -p bin
	go build -o bin/sheeit .

clean:
	rm -rf bin
