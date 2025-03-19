# Variables
APP_CONTROL := control
APP_SCANNER := scanner
APP_SHARED := shared
BUILD_DIR := build

# Default target
# all: build-control build-scanner
all: build-scanner

# Build CONTROL
build-control:
	@echo "Building CONTROL..."
	cd $(APP_CONTROL) && go build -o ../$(BUILD_DIR)/$(APP_CONTROL)

# Build Scanner
build-scanner:
	@echo "Building Scanner..."
	cd $(APP_SCANNER) && go build -o ../$(BUILD_DIR)/$(APP_SCANNER)

# # Run CONTROL
run-control:
	@echo "Running CONTROL..."
	cd $(APP_CONTROL) && go run $(APP_CONTROL).go

# Run Scanner
run-scanner:
	@echo "Running Scanner..."
	cd $(APP_SCANNER) && go run $(APP_SCANNER).go

# Tidy
tidy:
	@echo "Tidying..."
	cd $(APP_SCANNER) && go mod tidy
	cd $(APP_CONTROL) && go mod tidy
	cd $(APP_SHARED) && go mod tidy

# Clean up
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)