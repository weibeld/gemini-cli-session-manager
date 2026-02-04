.PHONY: build clean run testbed testbedrun

BINARY_NAME=geminictl
TESTBED_BIN=testbed
BUILD_DIR=bin

# Default testbed parameters
TESTBED_CONFIG=cmd/testbed/config/default.json
TESTBED_DIR=testbed-data

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/geminictl
	@echo "Building $(TESTBED_BIN)..."
	@go build -o $(BUILD_DIR)/$(TESTBED_BIN) ./cmd/testbed

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(TESTBED_DIR)/

run:
	@./$(BUILD_DIR)/$(BINARY_NAME) status $(ARGS)

testbed: build
	@echo "Generating test data in $(TESTBED_DIR)..."
	@./$(BUILD_DIR)/$(TESTBED_BIN) --config $(TESTBED_CONFIG) --dir $(TESTBED_DIR)

testbedrun: testbed
	@./$(BUILD_DIR)/$(BINARY_NAME) --testbed $(TESTBED_DIR) status
