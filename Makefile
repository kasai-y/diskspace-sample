BINARY_NAME=diskspace

all: build
run:
	go run main.go
build:
	go build -o bin/$(BINARY_NAME) -v .
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe
