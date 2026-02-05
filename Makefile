.PHONY: build build-app build-testbedgen clean run testbed run-testbed

APP_BIN=geminictl
TESTBEDGEN_BIN=testbedgen
BUILD_DIR=bin

# Default testbed parameters
DEFAULT_TESTBEDGEN_CONFIG=cmd/testbedgen/config/default.json
DEFAULT_TESTBED_DIR=tmp/testbed

build: build-app build-testbedgen

build-app:
	@echo "Building $(APP_BIN)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_BIN) ./cmd/geminictl

build-testbedgen:
	@echo "Building $(TESTBEDGEN_BIN)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(TESTBEDGEN_BIN) ./cmd/testbedgen

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp/

run: build-app
	@./$(BUILD_DIR)/$(APP_BIN) status $(ARGS)

testbed: build-testbedgen
	@./$(BUILD_DIR)/$(TESTBEDGEN_BIN) --config $(DEFAULT_TESTBEDGEN_CONFIG) --dir $(DEFAULT_TESTBED_DIR)

run-testbed: build-app
	@./$(BUILD_DIR)/$(APP_BIN) --testbed $(DEFAULT_TESTBED_DIR) status
