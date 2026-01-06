.PHONY: help test run build clean

help:
	@echo "Targets:"
	@echo "  test   Run tests"
	@echo "  run    Run server"
	@echo "  build  Build binary"
	@echo "  clean  Clean artifacts"

test:
	go test -v -count=1 ./...

run:
	RELAY_HTTP_ADDR=:8429 go run ./cmd/relay-ref

build:
	go build -o bin/relay-ref ./cmd/relay-ref

clean:
	go clean -testcache
	rm -rf bin
