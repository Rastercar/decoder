.PHONY: all get build test clean run run_dev 
.DEFAULT_GOAL: $(BIN_FILE)

PROJECT_NAME = reciever_ms

CMD_DIR = ./cmd

BIN_FILE = ./bin/$(PROJECT_NAME)

CONFIG_FILE = ./config/config.yml

DEV_CONFIG_FILE = ./config/config.dev.yml


# Get version constant
VERSION := $(shell git describe --abbrev=0 --tags --always)
BUILD := $(shell git rev-parse HEAD)

# Use linker flags to provide version/build settings to the binary
LDFLAGS=-ldflags "-s -w -X=main.version=$(VERSION) -X=main.build=$(BUILD)"

default: get build

# Fetch dependencies
get:
	@echo "[*] Downloading dependencies..."
	cd $(CMD_DIR) && go get
	@echo "[*] Finish..."

# Build the go binary
build:
	@echo "[*] Building $(PROJECT_NAME)..."
	go build $(LDFLAGS) -o $(BIN_FILE) $(CMD_DIR)/...
	@echo "[*] Finish..."

# Run all tests
test:
	go test -race -cover -coverprofile=coverage.out ./... 
	go tool cover -html=coverage.out

# Clears all compiled content
clean:
	rm -rf bin/
	rm -rf vendor/

# Builds and runs the aplication on production mode
run: build
	$(BIN_FILE) -config-file=$(CONFIG_FILE)

# Builds and runs the aplication on development mode
run_dev: build
	$(BIN_FILE) -config-file=$(DEV_CONFIG_FILE)