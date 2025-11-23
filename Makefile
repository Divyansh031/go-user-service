.PHONY: help proto build run test clean docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make proto        - Generate protobuf code"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make clean        - Clean generated files"

proto:
	@echo "Generating protobuf code..."
	@if [ ! -d "googleapis" ]; then \
		echo "Downloading googleapis..."; \
		rm -rf googleapis; \
		git clone --depth 1 https://github.com/googleapis/googleapis.git; \
	fi
	protoc \
		-I. \
		-I./googleapis \
		-I./api/proto \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		api/proto/user/v1/*.proto
	@echo "Protobuf code generated successfully!"

build: proto
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

run: build
	@echo "Running application..."
	./bin/server

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./tests/integration/...

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Waiting for ScyllaDB to be healthy...Please wait a minute."
	@until docker exec user-service-scylla cqlsh -u cassandra -p cassandra -e "describe keyspaces" >/dev/null 2>&1; do \
		echo "ScyllaDB not ready yet, waiting 5 seconds..."; \
		sleep 5; \
	done
	@echo "ScyllaDB is ready! Applying schema..."
	docker exec -i user-service-scylla cqlsh -u cassandra -p cassandra < scripts/schema.cql
	@echo "Database schema applied successfully!"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

clean:
	@echo "Cleaning generated files..."
	rm -rf bin/
	find . -name "*.pb.go" -delete
	find . -name "*.pb.gw.go" -delete
	rm -f coverage.out

.DEFAULT_GOAL := help