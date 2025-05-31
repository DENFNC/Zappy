APP_PROTO          = proto/grpc/v1
APP_PROTO_GEN      = proto/gen/v1
APP_PROTO_VALIDATE = proto/protoc-gen-validate
APP_PROTO_API      = proto/googleapis
APP_PROTO_HOOKS    = proto/grpc/v1/hooks
GOPATH            ?= $(shell go env GOPATH)
BIN_PATH          := $(GOPATH)/bin


ifndef POSTGRES_URL
$(error ❌ POSTGRES_URL is not set! Please provide it: make mig-up POSTGRES_URL=...)
endif


.PHONY: check-plugins mig-up mig-drop

mig-create:
	@echo "🚧 Creating migration files..."
	migrate create -ext sql -dir ./migrations init

mig-up:
	@echo "🚀 Migrating database with $(POSTGRES_URL)"
	migrate -source file://./migrations/ \
		-database $(POSTGRES_URL) up

mig-drop:
	@echo "💥 Dropping all database objects with $(POSTGRES_URL)"
	migrate -source file://./migrations/ \
		-database $(POSTGRES_URL) drop -f

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
