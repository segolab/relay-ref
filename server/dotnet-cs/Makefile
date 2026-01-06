.PHONY: help test run build clean

help:
	@echo "Targets:"
	@echo "  test   Run tests"
	@echo "  run    Run server"
	@echo "  build  Build server"
	@echo "  clean  Clean artifacts"

test:
	dotnet test relay-ref.sln

run:
	RELAY_HTTP_ADDR=:8429 dotnet run --project RelayRef.Server

build:
	dotnet build RelayRef.Server

clean:
	dotnet clean
