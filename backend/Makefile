# the output binary name
BINARY_NAME = build/bin/myapi
GO_MAIN_FILE = ./cmd/pokemonapi/main.go

GO_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# build the project
build: $(GO_FILES)
	@echo 'Building the project...'
	go build -o $(BINARY_NAME) $(GO_MAIN_FILE)

# Run the application
run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

# Format the code
fmt:
	@echo 'Formatting code...'
	go fmt ./...

lint:
	@echo 'Linting code...'
	go vet ./...

tidy:
	@echo "Tidying up dependencies..."
	go mod tidy

clean:
	@echo "Cleaning up..."
	rm -rf $(BINARY_NAME)

# Run all tasks
all: fmt lint tidy build

