# Variables
BINARY_NAME=mdns-browser
MAIN_PATH=cmd/mdns-browser/main.go
BUILD_DIR=bin

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	go test ./...

# Lint using golangci-lint in docker
.PHONY: lint
lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run

# Format code
.PHONY: format
format:
	goimports -w .

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Run the binary
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)