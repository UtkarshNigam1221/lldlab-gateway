.PHONY: build run run-local run-prod test clean deps fmt vet staticcheck lint check install-tools

build:
	go build -o gateway ./cmd/gateway

run: run-local

run-local:
	ENV=local go run ./cmd/gateway/main.go

run-prod:
	ENV=prod go run ./cmd/gateway/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	rm -f gateway coverage.out coverage.html
	go clean

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

lint:
	golangci-lint run

check: fmt vet staticcheck lint test
	@echo "✅ All checks passed!"

install-tools:
	@echo "Installing development tools..."
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Tools installed successfully"
