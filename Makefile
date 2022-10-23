.PHONY: deps deps-fmt deps-lint deps-test fmt fmt-check lint test

GOPATH ?= $(GOPATH)

deps: deps-fmt deps-lint deps-test

deps-fmt:
	go install mvdan.cc/gofumpt@latest

deps-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(GOPATH)/bin/mockgen:
	go install github.com/golang/mock/mockgen@latest

fmt:
	gofumpt -w .

fmt-check:
	gofumpt -d .

lint:
	golangci-lint run

generate: $(GOPATH)/bin/mockgen
	go generate ./...

test: generate
	go test -race ./...