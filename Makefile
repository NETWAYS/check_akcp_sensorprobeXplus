VERSION := $(shell git describe --tags --always)
BUILD := go build -v -ldflags "-s -w -X main.Version=$(VERSION)"

BINARY_BASE_NAME = check_akcp

.PHONY: all clean build test

BUILD_DIR = ./build

all: build test

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

build:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(BUILD) -o $(BUILD_DIR)/$(BINARY_BASE_NAME).amd64
	GOOS=windows GOARCH=amd64 $(BUILD) -o $(BUILD_DIR)/$(BINARY_BASE_NAME).amd64
