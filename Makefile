.PHONY: build test run clean
.PHONY: react-tui test-react typecheck-react verify-react-tui

BINARY=linuxlab

build:
	go build -o $(BINARY) ./cmd/linuxlab

test:
	go test ./... -v

react-tui:
	npm run tui:react

test-react:
	npm run test:react

typecheck-react:
	npm run typecheck:react

verify-react-tui:
	npm run test:react
	npm run typecheck-react
	go test ./... -v
	go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
