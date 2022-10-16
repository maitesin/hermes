.PHONY: deps fmt fmt-check lint test

deps:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

fmt:
	gofumpt -w .

fmt-check:
	gofumpt -d .

lint:
	golangci-lint run

test:
	go test -race ./...