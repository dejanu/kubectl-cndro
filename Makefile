BINARY := kubectl-cndro
CMD := ./cmd/$(BINARY)
DIST_DIR := dist
DARWIN_ARM64 := $(DIST_DIR)/$(BINARY)-darwin-arm64
LINUX_AMD64 := $(DIST_DIR)/$(BINARY)-linux-amd64

.PHONY: build install clean build-darwin-arm64 build-linux-amd64 build-cross

build:
	go build -o $(BINARY) $(CMD)

build-darwin-arm64:
	mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(DARWIN_ARM64) $(CMD)

build-linux-amd64:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_AMD64) $(CMD)

build-cross: build-darwin-arm64 build-linux-amd64

install:
	go install $(CMD)

clean:
	rm -f $(BINARY) $(DARWIN_ARM64) $(LINUX_AMD64)
