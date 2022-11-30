BINARY_NAME=happy-day-api
BINARY_OS=""
BINARY_ARCH=""

ifeq ($(OS),Windows_NT)
	BINARY_OS := windows
	ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
		BINARY_ARCH := amd64
	else
		ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
			BINARY_ARCH := amd64
		endif
	endif
else
	BINARY_OS := $(shell uname -s | tr A-Z a-z)
	UNAME_P = $(shell uname -m)
	ifeq ($(UNAME_P),x86_64)
		BINARY_ARCH := amd64
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		BINARY_ARCH := arm
	endif
endif

compile:
	GOOS=$(BINARY_OS) GOARCH=$(BINARY_ARCH) go build -o bin/$(BINARY_NAME) -ldflags "-s -w"

test:
	go test -v ./...

clean:
	GOOS=$(BINARY_OS) GOARCH=$(BINARY_ARCH) go clean
	rm -f bin/$(BINARY_NAME)