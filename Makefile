APP_NAME := go-jwt-auth
BUILD_DIR := build
BINARY := $(BUILD_DIR)/$(APP_NAME)

.PHONY: all build run clean test

all: build

build: 
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) ./cmd/api/main.go

run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)