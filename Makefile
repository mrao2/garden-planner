.PHONY: frontend-build backend-tidy up rebuild down logs

frontend-build:
	cd frontend && npm install && npm run build

backend-tidy:
	cd backend && go mod tidy

up:
	docker compose up -d

rebuild:
	$(MAKE) backend-tidy
	$(MAKE) frontend-build
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f --tail=200
