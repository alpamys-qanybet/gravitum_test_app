compose:
	docker compose -f docker-compose.yml up -d

test:
	go test -v ./...

run:
	go run cmd/gravitum-test-app/main.go

binary:
	go build ./cmd/gravitum-test-app

