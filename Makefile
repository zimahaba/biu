.DEFAULT_GOAL := build

build:
	go build -v -o ./ ./...

run:
	go run cmd/biu/biu.go