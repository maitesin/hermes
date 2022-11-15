.PHONY: fmt fmt-check lint generate test

$(GOPATH)/bin/gofumpt:
	go install mvdan.cc/gofumpt@latest

$(GOPATH)/bin/golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(GOPATH)/bin/mockgen:
	go install github.com/golang/mock/mockgen@latest

fmt: $(GOPATH)/bin/gofumpt
	gofumpt -w .

fmt-check: $(GOPATH)/bin/gofumpt
	gofumpt -d .

lint: generate $(GOPATH)/bin/golangci-lint
	golangci-lint run

generate: $(GOPATH)/bin/mockgen
	go generate ./...

test: generate
	go test -race ./...