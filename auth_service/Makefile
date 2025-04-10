APP_PROTO=proto/grpc/v1
APP_PROTO_GEN=proto/gen/v1
APP_PROTO_VALIDATE=proto/third_party
GOPATH ?= $(shell go env GOPATH)
BIN_PATH := $(GOPATH)/bin

.PHONY: grpc-gen grpc-clean grpc-rebuild


grpc-gen: check-plugins
	@echo "🔧 Generating gRPC code..."
	@mkdir -p $(APP_PROTO_GEN)
	PATH="$(BIN_PATH):$$PATH" \
	protoc -I=$(APP_PROTO) -I=$(APP_PROTO_VALIDATE) \
		--go_out=$(APP_PROTO_GEN) --go_opt=paths=source_relative \
		--go-grpc_out=$(APP_PROTO_GEN) --go-grpc_opt=paths=source_relative \
		--validate_out=$(APP_PROTO_GEN) \
		--validate_opt=lang=go \
		--validate_opt=paths=source_relative \
		$(wildcard $(APP_PROTO)/*.proto)


grpc-clean:
	@echo "🧹 Cleaning generated gRPC files..."
	rm -rf $(APP_PROTO_GEN)
	@echo "✅ Clean complete"


grpc-rebuild: grpc-clean grpc-gen


check-plugins:
	@echo "📦 Checking plugins in $(BIN_PATH)..."

	@if [ ! -x "$(BIN_PATH)/protoc-gen-go" ]; then \
		echo "⬇️  Installing protoc-gen-go..."; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
	else \
		echo "✅ protoc-gen-go found."; \
	fi

	@if [ ! -x "$(BIN_PATH)/protoc-gen-go-grpc" ]; then \
		echo "⬇️  Installing protoc-gen-go-grpc..."; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	else \
		echo "✅ protoc-gen-go-grpc found."; \
	fi

	@if [ ! -x "$(BIN_PATH)/protoc-gen-validate" ]; then \
		echo "⬇️  Installing protoc-gen-validate..."; \
		go get github.com/envoyproxy/protoc-gen-validate@latest; \
	else \
		echo "✅ protoc-gen-validate found."; \
	fi

	@if ! echo $$PATH | grep -q "$(BIN_PATH)"; then \
		echo "⚠️  $(BIN_PATH) is not in your PATH"; \
		echo "👉 Please add this to your shell config (e.g. ~/.bashrc):"; \
		echo "   export PATH=\"\$$PATH:$(BIN_PATH)\""; \
	fi