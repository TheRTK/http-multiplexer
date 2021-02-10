.PHONY: lint
lint:
	golangci-lint run --config ./build/ci/.golangci.yml ./...

.PHONY: build
build:
	go build -installsuffix 'static' -ldflags "-s -w" -o multiplexer cmd/multiplexer/*

