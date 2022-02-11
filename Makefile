.PHONY: build start buildLinux buildAndStart clean tool help

all: buildAndStart

build:
	@echo "build tjd_web to dist/"
	@go build -v -tags netgo -ldflags '-s -w' -o dist/tjd_web

start:
	./build/tjd_web

buildAndStart: build start

buildLinux:
	GOARCH=amd64 GOOS=linux go build -tags netgo -ldflags '-s -w' -o dist/tjd_linux_web

tool:
	go vet ./...; true
	gofmt -w .

clean:
	rm -rf dist/ var/
	go clean -i .

help:
	@echo "make build: compile packages and dependencies in local system"
	@echo "make buildLinux: compile packages for linux x86 system"
	@echo "make start: run ./build/tjd_web"
	@echo "make tool: run specified go tool"
	@echo "make clean: remove build files, object files and cached files"
