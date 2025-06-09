.PHONY: help build up down logs clean restart dev prod

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build all Docker images"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - Show logs from all services"
	@echo "  clean     - Remove all containers, images, and volumes"
	@echo "  restart   - Restart all services"
	@echo "  dev       - Start in development mode"
	@echo "  prod      - Start in production mode"
	@echo "  backend   - Start only backend services"
	@echo "  frontend  - Start only frontend service"
	@echo "  db        - Start only database service"

# Build all images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Clean everything
clean:
	docker-compose down -v --rmi all --remove-orphans
	docker system prune -f

# Restart services
restart: down up

# Development mode
dev:
	docker-compose up

# Production mode
prod:
	docker-compose up -d

# Start only backend services (database + backend)
backend:
	docker-compose up -d database backend

# Start only frontend
frontend:
	docker-compose up -d frontend

# Start only database
db:
	docker-compose up -d database

# Update and rebuild
update:
	git pull origin main
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Health check
health:
	@echo "Checking service health..."
	@docker-compose ps
	@echo "\nFrontend health:"
	@curl -f http://localhost:3000 > /dev/null 2>&1 && echo "✅ Frontend is healthy" || echo "❌ Frontend is not responding"
	@echo "Backend health:"
	@curl -f http://localhost:6050/health > /dev/null 2>&1 && echo "✅ Backend is healthy" || echo "❌ Backend is not responding"