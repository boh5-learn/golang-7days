lint:
	golangci-lint run ./...

tidy:
	go mod tidy

.PHONY: lint