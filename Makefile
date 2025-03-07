# Variables
APP_API := api
APP_SCANNER := scanner
APP_SHARED := shared
BUILD_DIR := build

# Default target
# all: build-api build-scanner
all: build-scanner

# Build API
build-api:
	@echo "Building API..."
	cd $(APP_API) && go build -o ../$(BUILD_DIR)/$(APP_API)

# Build Scanner
build-scanner:
	@echo "Building Scanner..."
	cd $(APP_SCANNER) && go build -o ../$(BUILD_DIR)/$(APP_SCANNER)

# # Run API
# run-api:
# 	@echo "Running API..."
# 	cd $(APP_API) && go run $(APP_API).go

# Run Scanner
run-scanner:
	@echo "Running Scanner..."
	cd $(APP_SCANNER) && go run $(APP_SCANNER).go

# Tidy
tidy:
	@echo "Tidying..."
	cd $(APP_SCANNER) && go mod tidy
# cd $(APP_API) && go mod tidy
# cd $(APP_SHARED) && go mod tidy

# Clean up
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)