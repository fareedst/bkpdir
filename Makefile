.PHONY: build-all build-ubuntu20 build-ubuntu22 build-ubuntu24 clean

# Build for all Ubuntu versions
build-all: build-ubuntu20 build-ubuntu22 build-ubuntu24

# Build for Ubuntu 20.04
build-ubuntu20:
	GOOS=linux GOARCH=amd64 go build -o bin/bkpdir-ubuntu20.04

# Build for Ubuntu 22.04
build-ubuntu22:
	GOOS=linux GOARCH=amd64 go build -o bin/bkpdir-ubuntu22.04

# Build for Ubuntu 24.04
build-ubuntu24:
	GOOS=linux GOARCH=amd64 go build -o bin/bkpdir-ubuntu24.04

# Clean build artifacts
clean:
	rm -rf bin/ 