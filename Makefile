BINARY_NAME=memrydb

all: build run

test:
	go test -v ./...

benchmark:
	go test -v -bench=. -benchmem -run=- ./...

build:
	go build -o ./bin/$(BINARY_NAME) ./cmd/server/main.go

run: build
	./bin/$(BINARY_NAME)

clean:
	rm -rf ./bin
	go clean

.PHONY: all test build run clean benchmark
