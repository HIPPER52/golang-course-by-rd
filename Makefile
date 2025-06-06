.PHONY: up down build clean

up:
	docker-compose up --build

down:
	docker-compose down

build:
	docker-compose build

clean:
	docker system prune -f
