migrate_up: 
	@go run ./cmd/migrate/main.go up
migrate_down: 
	@go run ./cmd/migrate/main.go down

build:
	@go build -o build/gin-event ./cmd/api
	@if [ -f ".env" ]; then \
		cp .env build/.env; \
	fi

exec:
	@nohup ./build/gin-event &
