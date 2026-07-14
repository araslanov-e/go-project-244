build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test -v ./...