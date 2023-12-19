.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/ddusage

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run
