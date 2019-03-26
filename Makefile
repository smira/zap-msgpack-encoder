all: test lint bench

.PHONY: test
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: bench
bench:
	go test -v -run=nothing -bench=. -benchmem .
