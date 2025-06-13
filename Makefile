BINARY_NAME=secm
BUILD_DIR=build
MAIN_PACKAGE=.
BUILD_NUM?=$(shell date +'%Y%m%d%H%M%S')

# Go build flags
VERSION?=$(shell git describe --tags)
LDFLAGS=-ldflags "-X github.com/open-zhy/secm/cmd.Version=${VERSION} -X github.com/open-zhy/secm/cmd.Build=${BUILD_NUM}"
GOFLAGS=-trimpath

# Example
# make install \
#	-e VERSION=$(git describe --tags || echo "$(git describe --tags $(git rev-list --tags --max-count=1))-$(git rev-parse --short HEAD)") \
#	-e BUILD_NUM="commit:$(git rev-parse HEAD)"

# Supported platforms
PLATFORMS=linux darwin windows
ARCHITECTURES=amd64 arm64

# Development tools
GOFMT=gofmt
GOTEST=go test

.PHONY: all build clean test fmt build-all build-plugins

# Default target
all: clean build

# Build for current platform
build: build-plugins
	@echo "Building ${BINARY_NAME}..."
	@go build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PACKAGE}

# Build for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@$(foreach PLATFORM,$(PLATFORMS),\
		$(foreach ARCH,$(ARCHITECTURES),\
			echo "Building for ${PLATFORM}/${ARCH}..." && \
			GOOS=${PLATFORM} GOARCH=${ARCH} go build ${GOFLAGS} ${LDFLAGS} \
				-o ${BUILD_DIR}/${BINARY_NAME}-${VERSION}-${PLATFORM}-${ARCH}$(if $(findstring windows,${PLATFORM}),.exe,) \
				${MAIN_PACKAGE} ; \
		)\
	)

# Build for specific platform (usage: make build-platform PLATFORM=darwin ARCH=arm64)
build-platform:
	@echo "Building for ${PLATFORM}/${ARCH}..."
	@GOOS=${PLATFORM} GOARCH=${ARCH} go build ${GOFLAGS} ${LDFLAGS} \
		-o ${BUILD_DIR}/${BINARY_NAME}-${VERSION}-${PLATFORM}-${ARCH}$(if $(findstring windows,${PLATFORM}),.exe,) \
		${MAIN_PACKAGE}

# Clean build directory
clean:
	@echo "Cleaning build directory..."
	@rm -rf ${BUILD_DIR}/*
	@mkdir -p ${BUILD_DIR}

# Run tests
test:
	@echo "Running tests..."
	@${GOTEST} -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@${GOFMT} -w .

# Install locally for development
install: build
	@echo "Installing ${BINARY_NAME}..."
	@cp ${BUILD_DIR}/${BINARY_NAME} ${GOPATH}/bin/

build-hello:
	@echo "Building Hello plugin..."
	@go build ${GOFLAGS} ${LDFLAGS} -buildmode=plugin -o ${BUILD_DIR}/plugins/hello.so plugins/hello/hello.go

build-p2p-transfer:
	@echo "Building \"P2P Transfer\" plugin..."
	@go build ${GOFLAGS} ${LDFLAGS} -buildmode=plugin -o ${BUILD_DIR}/plugins/p2p-transfer.so plugins/p2p-transfer/*.go

build-plugins: build-p2p-transfer

# Show help
help:
	@echo "Available targets:"
	@echo "  make              - Build for current platform"
	@echo "  make build-all    - Build for all platforms"
	@echo "  make build-platform PLATFORM=darwin ARCH=arm64  - Build for specific platform"
	@echo "  make clean        - Clean build directory"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format code"
	@echo "  make install      - Install locally"
	@echo ""
	@echo "Supported platforms: ${PLATFORMS}"
	@echo "Supported architectures: ${ARCHITECTURES}"
	@echo ""
	@echo "Example:"
	@echo "  make build-platform PLATFORM=darwin ARCH=arm64"
