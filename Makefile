.PHONY: all build clean install test run fmt vet

all: build

build:
	go build -o gen ./main.go

clean:
	go clean
	rm -f gen

install:
	go install ./...

test:
	go test ./... -v

run:
	go run ./main.go

fmt:
	go fmt ./...

vet:
	go vet ./...