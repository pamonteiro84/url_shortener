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

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l .

fmt-write:
	gofmt -w .

run:
	go run ./cmd/server
