all: generate test

generate:
	go generate ./...

test:
	go test ./...

