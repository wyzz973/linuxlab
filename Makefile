.PHONY: build test run clean

BINARY=linuxlab

build:
	go build -o $(BINARY) ./cmd/linuxlab

test:
	go test ./... -v

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
