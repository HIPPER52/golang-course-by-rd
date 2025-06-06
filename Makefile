.PHONY: up down build clean

up:
	docker-compose up --build

down:
	docker-compose down

build:
	docker-compose build

setup-envs:
	@echo "Setting up .env files for all services..."
	cp -n backend/.env.example backend/.env || echo "backend/.env already exists"
	cp -n admin-frontend/.env.example admin-frontend/.env || echo "admin-frontend/.env already exists"
	cp -n client-frontend/.env.example client-frontend/.env || echo "client-frontend/.env already exists"
	@echo "All .env files are set"

clean:
	docker system prune -f
