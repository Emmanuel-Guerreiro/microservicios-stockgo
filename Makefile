# Define variables
BINARY_NAME=stockgo
BIN_DIR=./bin
SRC_DIR=./

# Build command
build:
	mkdir -p $(BIN_DIR)               # Create bin directory if it doesn't exist
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR) # Compile the project

# Run command
run: build
	$(BIN_DIR)/$(BINARY_NAME)         # Execute the binary

# Clean command
clean:
	rm -rf $(BIN_DIR)                 # Remove the bin directory and all contents
