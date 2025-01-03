# note: call scripts from /scripts
SHELL:=/bin/bash
MAKE=make
BUILD_FLAGS = -ldflags="-s -w"
CURRENTDIR= $(shell pwd)
BIN = $(CURRENTDIR)/bin
BUILD_FLAGS = GO111MODULE=on GOOS=linux GOARCH=amd64
BUILD_TARGET = deployThis

.PHONY: zip run clean build

run-main:
	go run cmd/divergence/main.go
