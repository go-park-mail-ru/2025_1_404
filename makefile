test:
	go test ./... -cover

build:
	go build -o bin ./...

run:
	go run .