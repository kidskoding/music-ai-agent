BINARY_NAME=music_agent
ENTRY_POINT=./main.go

all: fmt vet build

build:
	go build -o $(BINARY_NAME) $(ENTRY_POINT)

run:
	go run $(ENTRY_POINT)

test:
	go test ./... -v

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)

tidy:
	go mod tidy
	go mod download

.PHONY: all build run test clean fmt vet tidy