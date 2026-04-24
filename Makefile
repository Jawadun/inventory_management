.PHONY: help build up down restart logs clean

help:
	@echo "Inventory Management System - Make Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make build     - Build all containers"
	@echo "  make up        - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make restart  - Restart all services"
	@echo "  make logs     - View logs"
	@echo "  make clean    - Stop and remove all volumes"
	@echo "  make db-reset - Reset database (WARNING: deletes data)"

build:
	docker-compose build

up:
	docker-compose up -d
	@echo ""
	@echo "Services started:"
	@echo "  Frontend:  http://localhost:3000"
	@echo "  Backend:   http://localhost:8080"
	@echo "  Database:  localhost:5432"

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

clean:
	docker-compose down -v

db-reset:
	docker-compose down -v
	@echo "Database has been reset. Run 'make up' to start fresh."

dev:
	@echo "Starting development environment..."
	docker-compose up -d postgres
	@echo ""
	@echo "PostgreSQL is running. Run migrations manually."
	@echo "To start backend: cd backend && go run ./cmd/server"
	@echo "To start frontend: cd frontend && npm run dev"