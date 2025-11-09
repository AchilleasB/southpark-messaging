# go-api-service

Purpose
- The Go HTTP API (producer) exposes `POST /messages`.
- Accepts JSON: `{ "author": "Cartman", "body": "Respect my authoritah!" }`.
- Uses Hexagonal Architecture:
  - internal/domain —> domain model and validation.
  - internal/ports —> Publisher interface (port).
  - internal/adapters/rabbitmq —> RabbitMQ adapter that implements Publisher.
  - internal/app —> application/service layer that marshals domain objects and calls the Publisher.
  - internal/handler —> HTTP handlers that adapt requests to the app.
  - cmd/api/main.go —> composition root (wires adapters, starts HTTP server).

Build & run (local)
```powershell
cd go-api-service
go mod tidy
go build ./main.go
# run locally
go run ./main.go
GO_API_PORT=8081 
```

Docker
- Image built via `go-api-service/Dockerfile`.
- Compose sets env to ensure `RABBITMQ_URL` is provided.