# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get


# Default target
all: test build

# Build the binary
build:
	$(GOBUILD) -o ./main cmd/server.go

# Test the project
test:
	$(GOTEST) -v ./...

# Clean the project
clean:
	$(GOCLEAN)

# Install dependencies
deps:
	$(GOGET) -v ./...

# Run the binary
run:
	./main

# Generate code coverage report
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Format the code using gofmt
fmt:
	gofmt -w .

# Vet the code for errors
vet:
	$(GOCMD) vet ./...

swagger:
	swag init -g internal/app/http.go  -o internal/docs

# Help target
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean up"
	@echo "  deps        - Install dependencies"
	@echo "  run         - Run the binary"
	@echo "  coverage    - Generate code coverage report"
	@echo "  fmt         - Format the code using gofmt"
	@echo "  vet         - Vet the code for errors"
	@echo "	 swagger	 - Generate swagger documentation"
	@echo "  help        - Show this help message"