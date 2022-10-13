.PHONY: deps fmt fmt-check lint

deps:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

fmt: deps
	gofumpt -w .

fmt-check: deps
	gofumpt -d .

lint: deps
	golangci-lint run

test:
	go test -race ./...