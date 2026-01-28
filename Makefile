BINARY=taps
BUILD_DIR=bin
SRC=./cmd/taps

.PHONY: build install clean run

build:
	go build -o $(BUILD_DIR)/$(BINARY) $(SRC)

install:
	go install $(SRC)

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY)
