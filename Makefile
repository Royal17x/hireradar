BINARY=hireradar
DC=docker compose

.PHONY: run build test coverage docker-up docker-down mock

run:
	go run cmd/bot/main.go
build:
	go build -o bin/$(BINARY) cmd/bot/main.go
test:
	go test ./internal/repository/postgres/... -v
	go test ./internal/usecase/... -v
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | grep total
docker-up:
	$(DC) up -d --build --remove-orphans
docker-down:
	$(DC) down -v --remove-orphans
mock:
	mockery --all --dir=internal/domain --output=internal/mocks --outpkg=mocks
