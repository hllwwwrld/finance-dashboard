# ============================================================================
# Конфигурация
# ============================================================================
MIGRATIONS_FOLDER:="./backend/db/migrations"
DB_DSN:="postgres://postgres:postgres@localhost:5432/dashboard?sslmode=disable"

bin-deps:
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOBIN=$(LOCAL_BIN) go install honnef.co/go/tools/cmd/staticcheck@latest

deps:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest


lint:
	golangci-lint run
	staticcheck ./backend/...

format:
	goimports -w ./backend/...

# ============================================================================
# Пересборка докер образа (медленная, без кеша)
# ============================================================================
# Использовать после изменения зависимостей или Dockerfile
build-clean:
	docker-compose build --no-cache backend frontend migrations

build-backend-clean:
	docker-compose build --no-cache backend

build-frontend-clean:
	docker-compose build --no-cache frontend



# ============================================================================
# Пересборка докер образа (быстрая, с кешем)
# ============================================================================
build-backend-fast:
	docker-compose build backend

build-frontend-fast:
	docker-compose build frontend

build-fast: build-backend-fast build-frontend-fast



# ============================================================================
# Локальная разработка (БЫСТРО для проверки интеграции!)
# ============================================================================
# Эти команды запускают сервисы локально без Docker (кроме БД)
# Идеально для быстрой проверки изменений и отладки

# Запуск только PostgreSQL в Docker
dev-db:
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 3
	goose -dir $(MIGRATIONS_FOLDER) postgres $(DB_DSN) up
	@echo "✓ Database is ready!"
	@echo "  Run 'make dev-backend' in another terminal to start backend"
	@echo "  Run 'make dev-frontend' in another terminal to start frontend"

# Запуск backend локально (требует запущенной БД: make dev-db)
dev-backend:
	@echo "Starting backend locally..."
	@echo "Make sure PostgreSQL is running (use 'make dev-db' to start it)"
	cd backend && go run cmd/main.go

# Запуск frontend локально (требует запущенного backend: make dev-backend)
dev-frontend:
	@echo "Starting frontend locally..."
	@echo "Make sure backend is running on http://localhost:3000"
	cd frontend && NEXT_PUBLIC_API_URL=http://localhost:3000 PORT=3001 pnpm dev

# Запуск всего стека локально в одном терминале
# ВАЖНО: Используйте отдельные терминалы для лучшего контроля логов
dev-all: dev-db
	@echo "Starting backend and frontend locally..."
	@echo "Backend: http://localhost:3000"
	@echo "Frontend: http://localhost:3001"
	@echo "Press Ctrl+C to stop all services"
	@trap 'kill 0' EXIT; \
	cd backend && go run cmd/main.go & \
	cd frontend && NEXT_PUBLIC_API_URL=http://localhost:3000 PORT=3001 pnpm dev & \
	wait




# Это надо разобрать, по хз хз чо сэтим делать, как-будто нужно только l-down
.PHONY: up
up:
	docker-compose up -d

build-up: build up

build-up-clean: build-clean up

.PHONY: down
down:
	docker-compose down



# Для локальной разработки
# с удалением volumes в контейнера постегреса, могут потеряться данные
l-up: up
	goose -dir $(MIGRATIONS_FOLDER) postgres $(DB_DSN) up

l-down:
	docker-compose down -v --remove-orphans

cleanup:
	docker system prune -a --volumes
