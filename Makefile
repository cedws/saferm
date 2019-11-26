BIN_DIR=bin
BIN_NAME=saferm.so
GOFLAGS=-trimpath

all: clean build

clean:
	rm -rf $(BIN_DIR)

build:
	go build -buildmode=c-shared $(GOFLAGS) -o $(BIN_DIR)/$(BIN_NAME) main.go
