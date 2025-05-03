APP_PROTO         = proto/grpc/v1
APP_PROTO_GEN     = proto/gen/v1
APP_PROTO_VALIDATE = proto/protoc-gen-validate
APP_PROTO_API     = proto/googleapis
GOPATH            ?= $(shell go env GOPATH)
BIN_PATH          := $(GOPATH)/bin

.PHONY: grpc-gen grpc-clean grpc-rebuild check-plugins mig-up mig-drop

# Generate gRPC and validation code
grpc-gen: check-plugins
	@echo "🔧 Generating gRPC code..."
	@mkdir -p $(APP_PROTO_GEN)
	PATH="$(BIN_PATH):$$PATH" \
	protoc \
		-I=$(APP_PROTO) \
		-I=$(APP_PROTO_VALIDATE) \
		-I=$(APP_PROTO_API) \
		--go_out=$(APP_PROTO_GEN) --go_opt=paths=source_relative \
		--go-grpc_out=$(APP_PROTO_GEN) --go-grpc_opt=paths=source_relative \
		--validate_out=$(APP_PROTO_GEN) \
		--validate_opt=lang=go \
		--validate_opt=paths=source_relative \
		--grpc-gateway_out=$(APP_PROTO_GEN) \
		--grpc-gateway_opt=paths=source_relative \
		$(wildcard $(APP_PROTO)/*.proto)

# Clean generated files
grpc-clean:
	@echo "🧹 Cleaning generated gRPC files..."
	rm -rf $(APP_PROTO_GEN)
	@echo "✅ Clean complete"

# Rebuild everything
grpc-rebuild: grpc-clean grpc-gen

mig-up:
	@echo "🚀 Migrating database..."
	migrate -source file://./migrations/ \
		-database pgx5://postgres:admin@localhost:5432/catalog_service up

mig-drop:
	@echo "💥 Dropping all database objects"
	migrate -source file://./migrations/ \
		-database pgx5://postgres:admin@localhost:5432/catalog_service drop -f

# Ensure required plugins and proto files are installed
check-plugins:
	@echo "📦 Checking plugins and proto definitions..."

	# Ensure proto-gen-validate definitions are present
	@if [ ! -d "$(APP_PROTO_VALIDATE)" ]; then \
		echo "⬇️  Cloning protoc-gen-validate definitions..."; \
		git clone https://github.com/bufbuild/protoc-gen-validate.git $(APP_PROTO_VALIDATE); \
	else \
		echo "✅ protoc-gen-validate definitions found."; \
		(cd $(APP_PROTO_VALIDATE) && git pull); \
	fi

	# Install protoc-gen-go if missing
	@if [ ! -x "$(BIN_PATH)/protoc-gen-go" ]; then \
		echo "⬇️  Installing protoc-gen-go binary..."; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
	else \
		echo "✅ protoc-gen-go found."; \
	fi

	# Install protoc-gen-go-grpc if missing
	@if [ ! -x "$(BIN_PATH)/protoc-gen-go-grpc" ]; then \
		echo "⬇️  Installing protoc-gen-go-grpc binary..."; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	else \
		echo "✅ protoc-gen-go-grpc found."; \
	fi

	# Build protoc-gen-validate binary if missing
	@if [ ! -x "$(BIN_PATH)/protoc-gen-validate" ]; then \
		echo "⬇️  Building protoc-gen-validate binary..."; \
		cd $(APP_PROTO_VALIDATE) && make build; \
	else \
		echo "✅ protoc-gen-validate binary found."; \
	fi

	# Warn if BIN_PATH not in PATH
	@if ! echo $$PATH | grep -q "$(BIN_PATH)"; then \
		echo "⚠️  $(BIN_PATH) is not in your PATH"; \
		echo "👉 Please add this to your shell config (e.g. ~/.bashrc):"; \
		echo "   export PATH='$$PATH:$(BIN_PATH)'"; \
	fi
