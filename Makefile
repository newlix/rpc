
# Run all tests.
test:
	@go test ./...
.PHONY: test

install:
	go install -v ./cmd/*
.PHONY: install