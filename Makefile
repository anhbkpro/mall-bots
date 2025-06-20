install-tools:
	@echo installing tools
	@go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@echo done

generate:
	@echo running code generation
	@go generate ./...
	@echo done

up:
	@echo "Starting monolith services..."
	docker-compose up --build -d
	@echo ""
	@echo "✅ Services started!"

	@echo "🔗 API: http://localhost:8080"

down:
	@echo "Stopping monolith services..."
	docker-compose down
	@echo ""
	@echo "✅ Services stopped!"

logs:
	@echo "Showing monolith logs..."
	docker-compose logs -f
	@echo ""
	@echo "✅ Logs shown!"
