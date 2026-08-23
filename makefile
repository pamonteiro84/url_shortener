docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f db
