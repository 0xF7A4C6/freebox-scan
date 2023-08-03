.PHONY: scan brute build

all:
	cls||clear
	@echo "Specify a target: scan, brute, build"

build:
	cls||clear
	cd ../cmd/freebruter && go build -o ../../scripts/

scan:
	cls||clear
	cd ../cmd/freebruter && go run . scan

brute:
	cls||clear
	cd ../cmd/freebruter && go run . brute
