lint:
	go mod tidy
	go fmt ./...
	go vet ./...
	golangci-lint run ./...
