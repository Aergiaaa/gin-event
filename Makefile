migrate_up: 
	@go run ./cmd/migrate/main.go up
migrate_down: 
	@go run ./cmd/migrate/main.go down

build:
	@go build -o build/gin-event ./cmd/api/main.go

exec:
	@nohup ./build/gin-event &