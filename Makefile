.PHONY: up down dev test-all logs

up:
	docker compose up --build -d

down:
	docker compose down

dev:
	docker compose -f docker-compose.dev.yml up --build

test-all:
	cd backend && go test ./...
	cd frontend && npm test -- --run

logs:
	docker compose logs -f
