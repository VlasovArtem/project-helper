build:
	go build -o bin/$(shell basename $(PWD)) cmd/main.go