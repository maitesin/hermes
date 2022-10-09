.PHONY: deps fmt fmt-check lint

deps:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

fmt:
	gofumpt -w .

fmt-check:
	gofumpt -d .

lint:
	golangci-lint run