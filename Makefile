APP_NAME := aegis-ai-gateway
DOCKER_COMPOSE := docker compose
DB_SERVICE := postgres
DB_USER := aegis
DB_NAME := aegis

.PHONY: help
help:
	@echo "Aegis AI Gateway commands:"
	@echo ""
	@echo "  make db-up       Start Postgres"
	@echo "  make db-down     Stop Docker services"
	@echo "  make db-logs     Show Postgres logs"
	@echo "  make db-shell    Open psql shell"
	@echo "  make api         Run API locally"
	@echo "  make worker      Run worker locally"
	@echo "  make test        Run Go tests"
	@echo "  make fmt         Format Go code"
	@echo "  make vet         Run go vet"
	@echo "  make tidy        Clean go.mod/go.sum"
	@echo "  make clean       Remove local build artifacts"


.PHONE: db-up
db-up:
	$(DOCKER_COMPOSE) up -d $(DB_SERVICE)

.PHONY: db-down
db-down:
	$(DOCKER_COMPOSE) down

.PHONY: db-logs
db-logs:
	$(DOCKER_COMPOSE) logs -f $(DB_SERVICE)

.PHONY: db-shell
db-shell:
	$(DOCKER_COMPOSE) exec $(DB_SERVICE) psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: api
api:
	go run ./cmd/api

.PHONY: worker
worker:
	go run ./cmd/worker

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
