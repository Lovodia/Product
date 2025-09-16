# Параметры
APP_NAME = product-app
MAIN_FILE = main.go
MIGRATIONS_DIR = ./migrations

# Переменные окружения 
YAML ?= local

# Команды

.PHONY: run
run:
	@echo "Launching application..."
	go run ./cmd/$(MAIN_FILE)

.PHONY: migrate
migrate:
	@echo "Starting migrations..."
	go run ./cmd/$(MAIN_FILE) --migrate

.PHONY: build
build:
	@echo "Assembly binary..."
	go build -o $(APP_NAME) $(MAIN_FILE)

.PHONY: clean
clean:
	@echo "Cleaning up collected file..."
	rm -f $(APP_NAME)

.PHONY: tidy
tidy:
	@echo "Updating dependencies"
	go mod tidy

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

.PHONY: test
test:
	@echo "Running unit tests with verbose output..."
	go test ./... -v

.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage report..."
	go test -cover ./internal/infrastructure/postgres

.PHONY: lint
lint:
	@echo "Running linters..."
	golangci-lint run ./... --timeout 3m