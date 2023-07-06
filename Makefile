MAIN_PATH := cmd/main.go
TOKEN_DIR := token
TOKEN_FILE_NAME := abitoff_read_adviser_bot.txt

build:
	go build $(MAIN_PATH)

run:
	mkdir -p $(TOKEN_DIR)
	go run $(MAIN_PATH) -tg-token="(shell cat $(TOKEN_DIR)/$(TOKEN_FILE_NAME) | tr -d '\n')"

clean:
	go clean