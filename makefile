colima-start:
	colima start

docker-up:
	colima start
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f db

build:
	go build ./...

run:
	go run ./cmd/server
