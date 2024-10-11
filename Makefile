.PHONY: lint
lint:
	golangci-lint run -v -c build/ci/.golangci-lint.yml

.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	go run cmd/detector/main.go
