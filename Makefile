.PHONY: release

BUILD_DIR := build
TARGET := sebuung-web

all: release

release: 
	rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@echo "Creating release for $(TARGET)..."
	# Add your release packaging commands here
	@echo "Release created successfully."
	@echo "Building binaries for all target platforms..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(TARGET)-linux-amd64
	GOOS=linux GOARCH=386 go build -o $(BUILD_DIR)/$(TARGET)-linux-386
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(TARGET)-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(TARGET)-windows-amd64.exe
	GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/$(TARGET)-windows-386.exe
	@echo "Binaries for all platforms created in $(BUILD_DIR)."