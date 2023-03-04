default: help

help:
	@echo "Makefile commands:"
	@echo ""
	@echo "- 'make build' - Build the project binary."
	@echo "- 'make update' - Update project dependencies."
	@echo "- 'make help' - Display list of available commands."

.PHONY: build
build:
	go build -o sumo

.PHONY: update
update:
	go get -u ./...
	go mod tidy -go=1.16 && go mod tidy -go=1.17

# https://stackoverflow.com/a/6273809/1826109
%:
	@:
