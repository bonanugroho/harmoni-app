APP_NAME=server
CMD_PATH=./cmd/server
BUILD_DIR=bin

.PHONY: build run test clean

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

test:
	go test -v -race -cover ./...

clean:
	rm -rf $(BUILD_DIR)
