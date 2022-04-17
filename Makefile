CURRENT_PATH = $(shell pwd)

go-lint: $(GOROOT) $(GOPATH)/bin/golangci-lint
	$(GOPATH)/bin/golangci-lint run --config=$(CURRENT_PATH)/.golangci.yaml --sort-results

entgen:
	go generate ./external/ent

.PHONY: setup
setup:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	# export PATH="$PATH:$(go env GOPATH)/bin"

.PHONY: protogen
protogen:
	protoc --go_out=. \
		--go-grpc_out=. \
		proto/*.proto

.PHONY: wiregen
wiregen:
	wire ./...

.PHONY: build
build:
	go build ./cmd/api

.PHONY: lint
lint:
	golangci-lint run