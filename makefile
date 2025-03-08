test:
	go test ./... -cover

build:
	go build -o bin/server ./...

run:
	go run .