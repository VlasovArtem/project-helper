build:
	go build -o bin/$(shell basename $(PWD)) cmd/main.go

unit-test:
	go test -v ./...

integration-test:
	go test -v ./... -tags=integration