BINARY_NAME=music_agent
MAIN_PKG=./cmd/agent

all: build

build:
	go build -o $(BINARY_NAME) $(MAIN_PKG)

run:
	go run $(MAIN_PKG)

clean:
	rm -f $(BINARY_NAME)

tidy:
	go mod tidy

.PHONY: test

test:
	go test ./... -v