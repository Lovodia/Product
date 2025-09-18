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

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

.PHONY: docker-up
docker-up:
	@echo "Starting docker-compose services..."
	docker-compose up -d --build

.PHONY: docker-down
docker-down:
	@echo "Stopping docker-compose services..."
	docker-compose down

.PHONY: docker-reset
docker-reset:
	@echo "Stopping and removing containers, volumes, networks..."
	docker-compose down -v

.PHONY: docker-logs
docker-logs:
	@echo "Showing logs for app service..."
	docker logs -f my_go_app

.PHONY: docker-migrate
docker-migrate:
	@echo "Running migrations inside Docker container..."
	docker exec -it my_go_app ./app --migrate

.PHONY: docker-seed
docker-seed:
	@echo "Seeding the database with dump_data.sql..."
	docker cp dump_data.sql my_postgres:/dump_data.sql
	docker exec -i my_postgres psql -U postgres -d practical3 -f /dump_data.sql

.PHONY: docker-reseed
docker-reseed:
	@echo "Restarting seed service..."
	docker-compose stop seed || true
	docker-compose rm -f seed || true
	docker-compose run --rm seed

.PHONY: docker-rebuild-run
docker-rebuild-run: docker-reset docker-build docker-up docker-seed
	@echo "Docker containers rebuilt, started and seeded."
