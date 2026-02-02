.PHONY: build clean run

BINARY_NAME=geminictl
BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/geminictl

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

run: build
	@./$(BUILD_DIR)/$(BINARY_NAME) status
