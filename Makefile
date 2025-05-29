default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build provider
.PHONY: build
build:
	go build -o terraform-provider-selectel-baremetal github.com/selectel/terraform-provider-selectel-baremetal/cmd/terraform-provider-selectel-baremetal

# Install provider locally for development
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/selectel/selectel-baremetal/0.1.0/darwin_arm64
	cp terraform-provider-selectel-baremetal ~/.terraform.d/plugins/registry.terraform.io/selectel/selectel-baremetal/0.1.0/darwin_arm64/

# Run unit tests
.PHONY: test
test:
	go test -v ./...

# Run unit tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Format code
.PHONY: fmt
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Generate documentation
.PHONY: docs
docs:
	go generate ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-selectel-baremetal
	rm -f coverage.out

# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Run all checks (format, lint, test)
.PHONY: check
check: fmt lint test

# Development setup
.PHONY: dev-setup
dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the provider binary"
	@echo "  install       - Install provider locally for development"
	@echo "  test          - Run unit tests"
	@echo "  testacc       - Run acceptance tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format Go and Terraform code"
	@echo "  lint          - Run linter"
	@echo "  docs          - Generate documentation"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  check         - Run format, lint, and test"
	@echo "  dev-setup     - Install development tools"
	@echo "  help          - Show this help message" 